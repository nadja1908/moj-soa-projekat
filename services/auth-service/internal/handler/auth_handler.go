package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"auth-service/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type AuthClaims struct {
	UserID int64  `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var (
	jwtKey        = []byte(os.Getenv("JWT_KEY"))
	refreshJwtKey = []byte(os.Getenv("JWT_REFRESH_KEY"))
)

type AuthHandler struct {
	stakeholdersServiceURL string
}

func NewAuthHandler(stakeholdersServiceURL string) *AuthHandler {
	if len(refreshJwtKey) == 0 {
		refreshJwtKey = append(jwtKey, []byte("-refresh")...)
	}
	return &AuthHandler{
		stakeholdersServiceURL: stakeholdersServiceURL,
	}
}

// Login handler - poziva stakeholders servis za autentifikaciju
func (h *AuthHandler) Login(c *gin.Context) {
	// Start tracing span
	tracer := otel.Tracer("auth-service")
	_, span := tracer.Start(c.Request.Context(), "login_user")
	defer span.End()

	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", "invalid_request"))
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"trace_id": span.SpanContext().TraceID().String(),
			"span_id":  span.SpanContext().SpanID().String(),
		}).Error("Invalid login request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(
		attribute.String("username", req.Username),
		attribute.String("operation", "user_login"),
	)

	logrus.WithFields(logrus.Fields{
		"username": req.Username,
		"trace_id": span.SpanContext().TraceID().String(),
		"span_id":  span.SpanContext().SpanID().String(),
	}).Info("User login attempt")

	// Pozovi stakeholders servis za login
	loginData, _ := json.Marshal(req)
	resp, err := http.Post(h.stakeholdersServiceURL+"/internal/login", "application/json", bytes.NewBuffer(loginData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to user service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		c.JSON(resp.StatusCode, errorResp)
		return
	}

	var userResp struct {
		User model.User `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
		return
	}

	// Generiši tokene
	accessToken, refreshToken, expiresAt, err := h.generateTokens(userResp.User.ID, userResp.User.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         userResp.User,
	})
}

// Register handler - poziva stakeholders servis za registraciju
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pozovi stakeholders servis za registraciju
	registerData, _ := json.Marshal(req)
	resp, err := http.Post(h.stakeholdersServiceURL+"/internal/register", "application/json", bytes.NewBuffer(registerData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to user service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		c.JSON(resp.StatusCode, errorResp)
		return
	}

	var userResp struct {
		User model.User `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
		return
	}

	// Generiši tokene
	accessToken, refreshToken, expiresAt, err := h.generateTokens(userResp.User.ID, userResp.User.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusCreated, model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         userResp.User,
	})
}

// VerifyToken verifikuje token
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	var req model.TokenValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := h.validateAccessToken(req.Token)
	if err != nil {
		c.JSON(http.StatusOK, model.TokenValidationResponse{
			Valid: false,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.TokenValidationResponse{
		Valid:  true,
		UserID: claims.UserID,
		Role:   claims.Role,
	})
}

// ValidateToken endpoint za validaciju tokena od strane drugih servisa
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.validateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId": claims.UserID,
		"role":   claims.Role,
	})
}

// RefreshToken generiše novi access token pomoću refresh token-a
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := h.validateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generiši nove tokene
	accessToken, refreshToken, expiresAt, err := h.generateTokens(claims.UserID, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"expiresAt":    expiresAt,
	})
}

// generateTokens generiše access i refresh tokene
func (h *AuthHandler) generateTokens(userID int64, role string) (string, string, time.Time, error) {
	// Access token (1 sat)
	expirationTime := time.Now().Add(1 * time.Hour)
	accessClaims := &AuthClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Refresh token (7 dana)
	refreshExpiration := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &AuthClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(refreshJwtKey)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessTokenString, refreshTokenString, expirationTime, nil
}

// GenerateJWT je public metoda za generiranje tokenai - koristi se u RPC handler
func (h *AuthHandler) GenerateJWT(userID int64, username, role string) (string, string, time.Time, error) {
	return h.generateTokens(userID, role)
}

// validateAccessToken validira access token
func (h *AuthHandler) validateAccessToken(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// validateRefreshToken validira refresh token
func (h *AuthHandler) validateRefreshToken(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshJwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	return claims, nil
}

// Health check
func (h *AuthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "auth-service",
	})
}
