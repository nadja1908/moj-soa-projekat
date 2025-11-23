package handler

import (
	"auth-service/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthRPCHandler implements the RPC interface for AuthHandler
type AuthRPCHandler struct {
	authHandler            *AuthHandler
	stakeholdersServiceURL string
}

func NewAuthRPCHandler(stakeholdersServiceURL string) *AuthRPCHandler {
	return &AuthRPCHandler{
		authHandler:            NewAuthHandler(stakeholdersServiceURL),
		stakeholdersServiceURL: stakeholdersServiceURL,
	}
}

// ProcessLogin handles login logic for RPC
func (h *AuthRPCHandler) ProcessLogin(req *model.LoginRequest) (*model.LoginResponse, error) {
	// Forward to stakeholders service
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(
		h.stakeholdersServiceURL+"/internal/login",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call stakeholders service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed")
	}

	var response struct {
		User struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Generate JWT tokens
	token, refreshToken, err := h.generateTokens(response.User.ID, response.User.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	return &model.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: model.UserInfo{
			ID:       response.User.ID,
			Username: response.User.Username,
			Email:    response.User.Email,
			Role:     response.User.Role,
		},
	}, nil
}

// ProcessRegister handles registration logic for RPC
func (h *AuthRPCHandler) ProcessRegister(req *model.RegisterRequest) (*model.RegisterResponse, error) {
	// Forward to stakeholders service
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(
		h.stakeholdersServiceURL+"/internal/register",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call stakeholders service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("registration failed")
	}

	var response struct {
		Message string `json:"message"`
		User    struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &model.RegisterResponse{
		Success: true,
		Message: "Registration successful",
		User: model.UserInfo{
			ID:       response.User.ID,
			Username: response.User.Username,
			Email:    response.User.Email,
			Role:     response.User.Role,
		},
	}, nil
}

// generateTokens helper function
func (h *AuthRPCHandler) generateTokens(userID int64, role string) (string, string, error) {
	jwtKey := []byte(os.Getenv("JWT_KEY"))
	refreshJwtKey := []byte(os.Getenv("JWT_REFRESH_KEY"))

	if len(refreshJwtKey) == 0 {
		refreshJwtKey = append(jwtKey, []byte("-refresh")...)
	}

	// Access token (15 minutes) - koristimo postojeću AuthClaims strukturu
	claims := &AuthClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Refresh token (7 days)
	refreshClaims := &AuthClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(refreshJwtKey)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}
