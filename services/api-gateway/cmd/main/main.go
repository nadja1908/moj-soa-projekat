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

	tourServiceURL := os.Getenv("TOUR_SERVICE_URL")
	log.Printf("DEBUG: Read TOUR_SERVICE_URL from env: '%s'", tourServiceURL)
	if tourServiceURL == "" {
		tourServiceURL = "http://tour-service:8004"
		log.Printf("DEBUG: Using default tour service URL: '%s'", tourServiceURL)
	}

	// Initialize Gin router
	router := gin.Default()

	// Disable automatic redirect for trailing slash to avoid CORS issues
	router.RedirectTrailingSlash = false

	// Configure CORS - include both old and new allowed origins
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
		"X-User-ID",
		"X-User-Role",
	}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	router.Use(cors.New(config))

	// Initialize handlers - include all services
	log.Printf("DEBUG: Creating gateway handler with tour URL: %s", tourServiceURL)
	gatewayHandler := handler.NewGatewayHandler(authServiceURL, stakeholdersServiceURL, blogServiceURL, tourServiceURL)
	log.Printf("DEBUG: Gateway handler created successfully!")

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
		auth.GET("/validate", gatewayHandler.ProxyToAuth)
	}

	// User routes (auth required)
	users := router.Group("/api/users")
	users.Use(authMiddleware.ValidateToken())
	{
		users.GET("/profile", gatewayHandler.ProxyToStakeholders)
		users.PUT("/profile", gatewayHandler.ProxyToStakeholders)
		users.DELETE("/profile", gatewayHandler.ProxyToStakeholders)
	}

	// Blog routes - blog-service handles its own auth
	blog := router.Group("/api/blog")
	{
		// Public routes
		blog.GET("/posts", gatewayHandler.ProxyToBlog)
		blog.GET("/posts/:id", gatewayHandler.ProxyToBlog)

		// Protected routes - blog-service has its own auth middleware
		blog.POST("/posts", gatewayHandler.ProxyToBlog)
		blog.PUT("/posts/:id", gatewayHandler.ProxyToBlog)
		blog.DELETE("/posts/:id", gatewayHandler.ProxyToBlog)
	}

	// Admin routes (auth required)
	admin := router.Group("/api/admin")
	admin.Use(authMiddleware.ValidateToken())
	{
		admin.GET("/users", gatewayHandler.ProxyToStakeholders)
		admin.PUT("/users/:id/block", gatewayHandler.ProxyToStakeholders)
		admin.PUT("/users/:id/unblock", gatewayHandler.ProxyToStakeholders)
		admin.GET("/posts", gatewayHandler.ProxyToBlog)
	}

	// Tour routes (some public, some auth required)
	log.Printf("DEBUG: About to configure tour routes...")
	tours := router.Group("/api/tours")
	{
		// Public routes
		tours.GET("/published", gatewayHandler.ProxyToTours)
		tours.GET("/public/:id", gatewayHandler.ProxyToTours)
	}

	// Protected tour routes (auth required)
	toursProtected := router.Group("/api/tours")
	toursProtected.Use(authMiddleware.ValidateToken())
	{
		// Tour CRUD (both with and without trailing slash)
		toursProtected.POST("", gatewayHandler.ProxyToTours)
		toursProtected.POST("/", gatewayHandler.ProxyToTours)
		toursProtected.GET("/my", gatewayHandler.ProxyToTours)
		toursProtected.PUT("/:id", gatewayHandler.ProxyToTours)
		toursProtected.DELETE("/:id", gatewayHandler.ProxyToTours)

		// Tour status management
		toursProtected.POST("/:id/publish", gatewayHandler.ProxyToTours)
		toursProtected.POST("/:id/archive", gatewayHandler.ProxyToTours)
		toursProtected.POST("/:id/reactivate", gatewayHandler.ProxyToTours)
	}
	log.Printf("DEBUG: Tour routes configured successfully!")

	// Key points routes (auth required)
	log.Printf("DEBUG: About to configure keypoints routes...")
	keypoints := router.Group("/api/keypoints")
	keypoints.Use(authMiddleware.ValidateToken())
	{
		keypoints.POST("", gatewayHandler.ProxyToTours)
		keypoints.GET("/:id", gatewayHandler.ProxyToTours)
		keypoints.PUT("/:id", gatewayHandler.ProxyToTours)
		keypoints.DELETE("/:id", gatewayHandler.ProxyToTours)
		keypoints.GET("/tour/:tourId", gatewayHandler.ProxyToTours)
		keypoints.POST("/reorder/:tourId", gatewayHandler.ProxyToTours)
	}
	log.Printf("DEBUG: Keypoints routes configured successfully!")

	log.Printf("API Gateway starting on port %s", port)
	log.Printf("Auth Service URL: %s", authServiceURL)
	log.Printf("Stakeholders Service URL: %s", stakeholdersServiceURL)
	log.Printf("Blog Service URL: %s", blogServiceURL)
	log.Printf("Tour Service URL: %s", tourServiceURL)
	log.Printf("Tour routes configured successfully!")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start API Gateway:", err)
	}
}
