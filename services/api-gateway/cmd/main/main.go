package main

import (
	"log"
	"os"

	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Get environment variables with defaults
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://auth-service:8003"
	}

	stakeholdersServiceURL := os.Getenv("STAKEHOLDERS_SERVICE_URL")
	if stakeholdersServiceURL == "" {
		stakeholdersServiceURL = "http://stakeholders-service:8001"
	}

	blogServiceURL := os.Getenv("BLOG_SERVICE_URL")
	if blogServiceURL == "" {
		blogServiceURL = "http://blog-service:8002"
	}

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Initialize handlers
	gatewayHandler := handler.NewGatewayHandler(authServiceURL, stakeholdersServiceURL, blogServiceURL)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authServiceURL)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "API Gateway"})
	})

	// Auth routes (no auth required)
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", gatewayHandler.ProxyToAuth)
		auth.POST("/register", gatewayHandler.ProxyToAuth)
		auth.POST("/refresh", gatewayHandler.ProxyToAuth)
	}

	// User routes (auth required)
	users := router.Group("/api/users")
	users.Use(authMiddleware.ValidateToken())
	{
		users.GET("/profile", gatewayHandler.ProxyToStakeholders)
		users.PUT("/profile", gatewayHandler.ProxyToStakeholders)
		users.DELETE("/profile", gatewayHandler.ProxyToStakeholders)
	}

	// Blog routes – AUTORIZACIJU radi SAMO blog-service
	blog := router.Group("/api/blog")
	{
		// javne rute
		blog.GET("/posts", gatewayHandler.ProxyToBlog)
		blog.GET("/posts/:id", gatewayHandler.ProxyToBlog)

		// zaštićene – ali gateway ih samo PROKSIRA,
		// a blog-service ima svoj AuthMiddleware koji proverava token
		blog.POST("/posts", gatewayHandler.ProxyToBlog)
		blog.PUT("/posts/:id", gatewayHandler.ProxyToBlog)
		blog.DELETE("/posts/:id", gatewayHandler.ProxyToBlog)
	}

	// Admin routes (auth required)
	admin := router.Group("/api/admin")
	admin.Use(authMiddleware.ValidateToken())
	{
		admin.GET("/users", gatewayHandler.ProxyToStakeholders)
		admin.DELETE("/users/:id", gatewayHandler.ProxyToStakeholders)
		admin.GET("/posts", gatewayHandler.ProxyToBlog)
	}

	log.Printf("API Gateway starting on port %s", port)
	log.Printf("Auth Service URL: %s", authServiceURL)
	log.Printf("Stakeholders Service URL: %s", stakeholdersServiceURL)
	log.Printf("Blog Service URL: %s", blogServiceURL)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start API Gateway:", err)
	}
}
