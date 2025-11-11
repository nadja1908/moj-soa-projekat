package middleware

import (
	"encoding/json"
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
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token (remove "Bearer " prefix)
		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Make request to auth service with token in Authorization header
		req, err := http.NewRequest("GET", m.authServiceURL+"/validate", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare validation request"})
			c.Abort()
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
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

		c.Next()
	}
}