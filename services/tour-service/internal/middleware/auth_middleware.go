package middleware

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware čita X-User-ID header i postavlja userId u context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDHeader := c.GetHeader("X-User-ID")
		userRole := c.GetHeader("X-User-Role")

		log.Printf("DEBUG: AuthMiddleware - X-User-ID header: '%s', X-User-Role: '%s'", userIDHeader, userRole)

		if userIDHeader != "" {
			userID, err := strconv.ParseInt(userIDHeader, 10, 64)
			if err == nil {
				c.Set("userId", userID)
				c.Set("userRole", userRole)
				log.Printf("DEBUG: AuthMiddleware - Set userId: %d, userRole: %s", userID, userRole)
			} else {
				log.Printf("DEBUG: AuthMiddleware - Failed to parse userID: %v", err)
			}
		} else {
			log.Printf("DEBUG: AuthMiddleware - No X-User-ID header found")
		}

		c.Next()
	}
}
