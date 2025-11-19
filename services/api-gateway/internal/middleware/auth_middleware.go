package middleware

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authServiceURL string
}

func NewAuthMiddleware(authServiceURL string) *AuthMiddleware {
	return &AuthMiddleware{
		authServiceURL: authServiceURL,
	}
}

func (m *AuthMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("DEBUG: ValidateToken STARTED for path: %s", c.Request.URL.Path)

		// Get Authorization header from original request (browser -> gateway)
		authHeader := c.GetHeader("Authorization")
		log.Printf("DEBUG: ValidateToken - Auth header: '%s'", authHeader)

		if authHeader == "" {
			log.Println("[AuthMiddleware] Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("[AuthMiddleware] Invalid Authorization header format: %s\n", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract token (remove "Bearer " prefix)
		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		if token == "" {
			log.Println("[AuthMiddleware] Empty token after Bearer prefix")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Make request to auth service for validation
		validateURL := m.authServiceURL + "/validate"
		log.Printf("[AuthMiddleware] Validating token at %s\n", validateURL)

		req, err := http.NewRequest(http.MethodGet, validateURL, nil)
		if err != nil {
			log.Printf("[AuthMiddleware] Failed to create validation request: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create validation request"})
			c.Abort()
			return
		}

		// Forward the entire Authorization header to auth service
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[AuthMiddleware] Failed to call auth-service: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Read response body for debug
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := string(bodyBytes)

		if resp.StatusCode != http.StatusOK {
			log.Printf(
				"[AuthMiddleware] Auth-service rejected token. Status=%d, Body=%s\n",
				resp.StatusCode,
				bodyStr,
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Parse JSON response from auth service
		var validationResponse struct {
			UserID   int64  `json:"userId"`
			Username string `json:"username"`
			Role     string `json:"role"`
		}

		if err := json.Unmarshal(bodyBytes, &validationResponse); err != nil {
			log.Printf("[AuthMiddleware] Failed to parse auth response: %v, body=%s\n", err, bodyStr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse validation response"})
			c.Abort()
			return
		}

		log.Printf(
			"[AuthMiddleware] Token valid. userId=%d, username=%s, role=%s\n",
			validationResponse.UserID,
			validationResponse.Username,
			validationResponse.Role,
		)

		// Store user data in context for downstream services
		c.Set("userID", int(validationResponse.UserID)) // Convert to int for compatibility
		c.Set("username", validationResponse.Username)
		c.Set("userRole", validationResponse.Role) // Note: using "userRole" for compatibility with tour service
		c.Set("role", validationResponse.Role)

		log.Printf("DEBUG: ValidateToken SUCCESS - Set userID: %d, userRole: %s", int(validationResponse.UserID), validationResponse.Role)

		// Continue to next handler
		c.Next()
	}
}
