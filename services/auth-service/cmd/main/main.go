package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"auth-service/internal/handler"
)

func main() {
	port := getEnv("PORT", "8003")
	stakeholdersService := getEnv("STAKEHOLDERS_SERVICE_URL", "http://stakeholders-service:8001")

	// Inicijalizacija auth handler-a
	authHandler := handler.NewAuthHandler(stakeholdersService)

	r := gin.Default()

	// CORS konfiguracija
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", authHandler.Health)

	// Auth endpoints
	r.POST("/login", authHandler.Login)
	r.POST("/register", authHandler.Register)
	r.POST("/verify", authHandler.VerifyToken)
	r.POST("/refresh", authHandler.RefreshToken)

	// Token validation endpoint za druge servise
	r.GET("/validate", authHandler.ValidateToken)

	log.Printf("Auth service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}