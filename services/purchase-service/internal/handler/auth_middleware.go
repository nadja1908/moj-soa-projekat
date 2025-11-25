package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://auth-service:8003"
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		req, err := http.NewRequest("GET", authServiceURL+"/validate", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create validation request"})
			c.Abort()
			return
		}

		req.Header.Set("Authorization", authHeader)

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

		var authResp struct {
			UserID int64  `json:"userId"`
			Role   string `json:"role"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse auth response"})
			c.Abort()
			return
		}

		c.Set("userID", authResp.UserID)
		c.Set("userRole", authResp.Role)
		c.Next()
	}
}
