package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		fmt.Printf("========= MIDDLEWARE STARTED for path: %s ========\n", c.Request.URL.Path)
		log.Printf("DEBUG: ValidateToken STARTED for path: %s", c.Request.URL.Path)

		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("========= AUTH HEADER: '%s' ========\n", authHeader)
		log.Printf("DEBUG: ValidateToken - Auth header: '%s'", authHeader)

		if authHeader == "" {
			log.Printf("DEBUG: ValidateToken - No auth header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token (remove "Bearer " prefix)
		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			log.Printf("DEBUG: ValidateToken - Invalid auth header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		log.Printf("DEBUG: ValidateToken - Extracted token: '%s'", token[:20]+"...")

		// Make request to auth service with token in Authorization header
		validateURL := m.authServiceURL + "/validate"
		req, err := http.NewRequest("GET", validateURL, nil)
		if err != nil {
			log.Printf("DEBUG: ValidateToken - Failed to create request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare validation request"})
			c.Abort()
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)

		log.Printf("DEBUG: ValidateToken - Calling auth service: %s", validateURL)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("DEBUG: ValidateToken - Auth service request failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		log.Printf("DEBUG: ValidateToken - Auth service response status: %d", resp.StatusCode)

		if resp.StatusCode != http.StatusOK {
			log.Printf("DEBUG: ValidateToken - Token validation failed")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Parse response to get user info
		var validationResponse struct {
			UserID int    `json:"userId"`
			Role   string `json:"role"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&validationResponse); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse validation response"})
			c.Abort()
			return
		}

		// Set user info in context for stakeholders service
		c.Set("userID", validationResponse.UserID)
		c.Set("userRole", validationResponse.Role)

		fmt.Printf("========= MIDDLEWARE SUCCESS - Set userID: %d, userRole: %s ========\n", validationResponse.UserID, validationResponse.Role)
		log.Printf("DEBUG: ValidateToken SUCCESS - Set userID: %d, userRole: %s", validationResponse.UserID, validationResponse.Role)

		c.Next()
	}
}
