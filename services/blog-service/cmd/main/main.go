package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"blog-service/internal/handler"
	"blog-service/internal/store"
)

func main() {
	// Čitanje environment varijabli
	port := getEnv("PORT", "8002")
	dbUser := getEnv("DB_USER", "user")
	dbPass := getEnv("DB_PASS", "password")
	dbHost := getEnv("DB_HOST", "localhost")
	dbName := getEnv("DB_NAME", "blog_db")

	// Inicijalizacija store-a
	store, err := store.NewStore(dbUser, dbPass, dbHost, dbName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer store.Close()

	// Inicijalizacija handler-a
	blogHandler := handler.NewBlogHandler(store)

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
	r.GET("/health", blogHandler.Health)

	// Public routes
	r.GET("/posts", blogHandler.GetAllBlogPosts)
	r.GET("/posts/:id", blogHandler.GetBlogPost)

	// Protected routes
	protected := r.Group("/")
	protected.Use(handler.AuthMiddleware())
	{
		// Blog post endpoints
		protected.POST("/posts", blogHandler.CreateBlogPost)
		protected.POST("/posts/:id/comments", blogHandler.CreateComment)
		protected.POST("/posts/:id/like", blogHandler.LikeBlogPost)
		protected.DELETE("/posts/:id/like", blogHandler.UnlikeBlogPost)
	}

	log.Printf("Blog service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}