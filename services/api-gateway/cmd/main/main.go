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

	purchaseServiceURL := os.Getenv("PURCHASE_SERVICE_URL")
	log.Printf("DEBUG: Read PURCHASE_SERVICE_URL from env: '%s'", purchaseServiceURL)
	if purchaseServiceURL == "" {
		purchaseServiceURL = "http://purchase-service:8005"
		log.Printf("DEBUG: Using default purchase service URL: '%s'", purchaseServiceURL)
	}

	followerServiceURL := os.Getenv("FOLLOWER_SERVICE_URL")
	if followerServiceURL == "" {
		followerServiceURL = "http://follower-service:8006"
	}

	// RPC addresses
	authRPCAddr := "auth-service:9003"
	tourRPCAddr := "tour-service:9004"

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
	gatewayHandler := handler.NewGatewayHandler(authServiceURL, stakeholdersServiceURL, blogServiceURL, tourServiceURL, purchaseServiceURL)
	log.Printf("DEBUG: Gateway handler created successfully!")

	// Initialize RPC handler
	rpcHandler, err := handler.NewRPCHandler(authRPCAddr, tourRPCAddr, gatewayHandler)
	if err != nil {
		log.Printf("WARNING: Failed to initialize RPC handler: %v", err)
		log.Printf("RPC endpoints will not be available")
	} else {
		defer rpcHandler.Close()
		log.Printf("RPC handler initialized successfully!")
	}

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

	// RPC Auth routes (alternative implementation)
	if rpcHandler != nil {
		rpcAuth := router.Group("/api/rpc/auth")
		{
			rpcAuth.POST("/login", rpcHandler.LoginRPC)
			rpcAuth.POST("/register", rpcHandler.RegisterRPC)
		}
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
		// Public
		blog.GET("/posts", gatewayHandler.ProxyToBlog)
		blog.GET("/posts/:id", gatewayHandler.ProxyToBlog)

		// Protected
		blog.POST("/posts", gatewayHandler.ProxyToBlog)
		blog.PUT("/posts/:id", gatewayHandler.ProxyToBlog)
		blog.DELETE("/posts/:id", gatewayHandler.ProxyToBlog)

		// Upload slika
		blog.POST("/uploads", gatewayHandler.ProxyToBlog)
		blog.GET("/uploads/*filepath", gatewayHandler.ProxyToBlog)
	}

	log.Printf("DEBUG: About to configure purchase/cart routes...")
	cart := router.Group("/api/purchase")
	cart.Use(authMiddleware.ValidateToken())
	{
		cart.GET("", gatewayHandler.ProxyToPurchase)
		cart.GET("/", gatewayHandler.ProxyToPurchase)
		cart.POST("/add", gatewayHandler.ProxyToPurchase)
		cart.DELETE("/:tourId", gatewayHandler.ProxyToPurchase)
		cart.POST("/checkout", gatewayHandler.ProxyToPurchase)
	}
	log.Printf("DEBUG: Purchase routes configured successfully!")

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
		// Public routes - standard HTTP proxy
		tours.GET("/published", gatewayHandler.ProxyToTours)
		tours.GET("/public/:id", gatewayHandler.ProxyToTours)
	}

	// RPC Tour routes - proxy to HTTP endpoints that use RPC internally
	rpcTours := router.Group("/api/rpc/tours")
	{
		rpcTours.GET("/published", func(c *gin.Context) {
			// Proxy to tour service RPC endpoint
			gatewayHandler.ProxyRequest(c, gatewayHandler.TourServiceURL()+"/api/rpc/tours/published")
		})
		rpcTours.GET("/public/:id", func(c *gin.Context) {
			// Proxy to tour service RPC endpoint
			id := c.Param("id")
			gatewayHandler.ProxyRequest(c, gatewayHandler.TourServiceURL()+"/api/rpc/tours/public/"+id)
		})
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

	// RPC Tour routes (alternative implementation for specific use cases)
	if rpcHandler != nil {
		rpcToursAuth := router.Group("/api/rpc/tours")
		rpcToursAuth.Use(authMiddleware.ValidateToken())
		{
			rpcToursAuth.POST("", rpcHandler.RPCCreateTour)
			rpcToursAuth.POST("/", rpcHandler.RPCCreateTour)
			rpcToursAuth.GET("/my", rpcHandler.RPCGetTours)
		}
	}
	log.Printf("DEBUG: Tour routes configured successfully!")

	// Key points routes (auth required)
	log.Printf("DEBUG: About to configure keypoints routes...")
	keypoints := router.Group("/api/keypoints")
	keypoints.Use(authMiddleware.ValidateToken())
	{
		keypoints.POST("", gatewayHandler.ProxyToKeyPoints)
		keypoints.GET("/:id", gatewayHandler.ProxyToKeyPoints)
		keypoints.PUT("/:id", gatewayHandler.ProxyToKeyPoints)
		keypoints.DELETE("/:id", gatewayHandler.ProxyToKeyPoints)
		keypoints.GET("/tour/:tourId", gatewayHandler.ProxyToKeyPoints)
		keypoints.POST("/reorder/:tourId", gatewayHandler.ProxyToKeyPoints)
	}
	log.Printf("DEBUG: Keypoints routes configured successfully!")

	// Review routes (Public and Protected)
	log.Printf("DEBUG: About to configure review routes...")

	reviews := router.Group("/api/reviews")
	{
		// Public GET
		reviews.GET("/tour/:tourId", gatewayHandler.ProxyToTours)
	}

	reviewsProtected := router.Group("/api/reviews")
	reviewsProtected.Use(authMiddleware.ValidateToken())
	{
		// Protected POST
		reviewsProtected.POST("", gatewayHandler.ProxyToTours)
		reviewsProtected.POST("/", gatewayHandler.ProxyToTours)
	}

	log.Printf("DEBUG: Review routes configured successfully!")

	log.Printf("API Gateway starting on port %s", port)
	log.Printf("Auth Service URL: %s", authServiceURL)
	log.Printf("Stakeholders Service URL: %s", stakeholdersServiceURL)
	log.Printf("Blog Service URL: %s", blogServiceURL)
	log.Printf("Tour Service URL: %s", tourServiceURL)
	log.Printf("Tour routes configured successfully!")
	log.Printf("Purchase Service URL: %s", purchaseServiceURL)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start API Gateway:", err)
	}
}
