package main

import (
	"log"
	"net/http"
	"os"

	"blog-service/internal/handler"
	"blog-service/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Build version 4 - comments count added
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
	stakeholdersServiceURL := getEnv("STAKEHOLDERS_SERVICE_URL", "http://stakeholders-service:8001")
	followerServiceURL := getEnv("FOLLOWER_SERVICE_URL", "http://follower-service:8080")
	blogHandler := handler.NewBlogHandler(store, stakeholdersServiceURL, followerServiceURL)

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// 👉 Static serving za slike
	r.Static("/uploads", "./uploads")

	// Health check endpoint
	r.GET("/health", blogHandler.Health)

	// Public routes (sa opcional auth za isLiked)
	r.GET("/posts", handler.OptionalAuthMiddleware(), blogHandler.GetAllBlogPosts)
	r.GET("/posts/:id", blogHandler.GetBlogPost)

	// Protected routes
	protected := r.Group("/")
	protected.Use(handler.AuthMiddleware())
	{
		protected.POST("/posts", blogHandler.CreateBlogPost)
		protected.POST("/posts/:id/comments", blogHandler.CreateComment)
		protected.PUT("/comments/:commentId", blogHandler.UpdateComment)
		protected.DELETE("/comments/:commentId", blogHandler.DeleteComment)
		protected.POST("/posts/:id/like", blogHandler.LikeBlogPost)
		protected.DELETE("/posts/:id/like", blogHandler.UnlikeBlogPost)

		// 👉 UPLOAD SLIKA
		protected.POST("/uploads", blogHandler.UploadImage)
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
