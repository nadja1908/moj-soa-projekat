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
		// 1) Uzmemo Authorization header iz ORIGINALNOG zahteva (browser -> gateway)
		authHeader := c.GetHeader("Authorization")
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

		// Sam token (bez "Bearer ")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		if token == "" {
			log.Println("[AuthMiddleware] Empty token after Bearer prefix")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 2) Napravimo GET /validate ka AUTH servisu – ISTO kao što radi blog-service
		validateURL := m.authServiceURL + "/validate"
		log.Printf("[AuthMiddleware] Validating token at %s\n", validateURL)

		req, err := http.NewRequest(http.MethodGet, validateURL, nil)
		if err != nil {
			log.Printf("[AuthMiddleware] Failed to create validation request: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create validation request"})
			c.Abort()
			return
		}

		// Prosledimo ceo Authorization header (Bearer + token)
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

		// Pročitamo telo radi debug-a
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

		// 3) Parsiranje JSON odgovora auth-servisa
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

		// 4) Snimimo user podatke u context da ih dalje servisi mogu koristiti ako hoćeš
		c.Set("userID", validationResponse.UserID)
		c.Set("username", validationResponse.Username)
		c.Set("role", validationResponse.Role)

		// Nastavi dalje
		c.Next()
	}
}
