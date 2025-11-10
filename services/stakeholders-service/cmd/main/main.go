package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"stakeholders-service/internal/handler"
	"stakeholders-service/internal/store"
)

func main() {
	// Čitanje environment varijabli
	port := getEnv("PORT", "8001")
	dbUser := getEnv("DB_USER", "user")
	dbPass := getEnv("DB_PASS", "password")
	dbHost := getEnv("DB_HOST", "localhost")
	dbName := getEnv("DB_NAME", "stakeholders_db")

	// Inicijalizacija store-a (konekcija sa bazom)
	store, err := store.NewStore(dbUser, dbPass, dbHost, dbName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer store.Close()

	// Inicijalizacija handler-a
	userHandler := handler.NewUserHandler(store)

	// Gin setup
	r := gin.Default()

	// CORS konfiguracija
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	r.GET("/health", userHandler.Health)

	// Internal routes (za Auth servis)
	internal := r.Group("/internal")
	{
		internal.POST("/login", userHandler.Login)
		internal.POST("/register", userHandler.Register)
		internal.GET("/users/:id", userHandler.GetUserByID)
	}

	// Protected routes (sa autentifikacijom kroz Auth servis)
	protected := r.Group("/")
	protected.Use(handler.AuthMiddleware())
	{
		// Administrator endpoints
		protected.GET("/users", userHandler.GetAllUsers)
		protected.PUT("/users/:id/block", userHandler.BlockUser)
	}

	log.Printf("Stakeholders service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}