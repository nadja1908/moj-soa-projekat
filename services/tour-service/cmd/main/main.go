package main

import (
	"log"
	"os"
	"strconv"

	"tour-service/internal/handler"
	"tour-service/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Environment varijable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8004"
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" {
		log.Fatal("Database environment variables must be set")
	}

	// Inicijalizacija store
	store := store.NewStore(dbUser, dbPass, dbHost, dbName)
	defer store.Close()

	// Inicijalizacija handlers
	tourHandler := handler.NewTourHandler(store)
	keyPointHandler := handler.NewKeyPointHandler(store)
	durationHandler := handler.NewTourDurationHandler(store)
	simulatorHandler := handler.NewTourSimulatorHandler(store)

	// Gin setup
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// JWT middleware function - sada čita user info iz headers umesto parsiranja JWT
	authMiddleware := func(c *gin.Context) {
		// Čitaj user info iz headers koje šalje API Gateway
		userIDHeader := c.GetHeader("X-User-ID")
		userRoleHeader := c.GetHeader("X-User-Role")

		log.Printf("DEBUG: Tour authMiddleware - X-User-ID: '%s', X-User-Role: '%s'", userIDHeader, userRoleHeader)

		if userIDHeader == "" {
			log.Printf("DEBUG: Tour authMiddleware - No X-User-ID header!")
			c.JSON(401, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Set userId in context for handlers
		userID, err := strconv.ParseInt(userIDHeader, 10, 64)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		c.Set("userId", userID)
		log.Printf("DEBUG: Tour authMiddleware - Calling c.Next()")
		c.Next()
		log.Printf("DEBUG: Tour authMiddleware - After c.Next()")
	}

	// Static files za slike
	r.Static("/uploads", "./uploads")

	// Health check
	r.GET("/health", tourHandler.Health)

	// Public rute (bez autentifikacije) - tours
	r.GET("/api/tours/published", tourHandler.GetPublishedTours)
	r.GET("/api/tours/public/:id", tourHandler.GetTourForTourist) // For tourists - only first key point

	// Public rute - key points i durations (koristim drugačiji path)
	r.GET("/api/keypoints/tour/:tourId", keyPointHandler.GetTourKeyPoints)
	r.GET("/api/durations/tour/:tourId", durationHandler.GetTourDurations)

	// Protected tour routes (defined individually to avoid group slash issues)
	r.POST("/api/tours", authMiddleware, tourHandler.CreateTour)
	r.GET("/api/tours/my", authMiddleware, tourHandler.GetMyTours)
	r.PUT("/api/tours/:id", authMiddleware, tourHandler.UpdateTour)
	r.POST("/api/tours/:id/publish", authMiddleware, tourHandler.PublishTour)
	r.POST("/api/tours/:id/archive", authMiddleware, tourHandler.ArchiveTour)
	r.POST("/api/tours/:id/reactivate", authMiddleware, tourHandler.ReactivateTour)
	r.DELETE("/api/tours/:id", authMiddleware, tourHandler.DeleteTour)

	// Key points management (original routes)
	r.POST("/api/tours/keypoints", authMiddleware, keyPointHandler.CreateKeyPoint)
	r.GET("/api/tours/keypoints/:id", authMiddleware, keyPointHandler.GetKeyPoint)
	r.PUT("/api/tours/keypoints/:id", authMiddleware, keyPointHandler.UpdateKeyPoint)
	r.DELETE("/api/tours/keypoints/:id", authMiddleware, keyPointHandler.DeleteKeyPoint)

	// Key points management (Gateway-compatible routes)
	// Gateway maps: POST /api/keypoints -> POST /api/tours/keypoints (za kreiranje)
	// Gateway maps: GET /api/keypoints/tour/123 -> GET /api/tours/tour/123 (za čitanje)
	r.GET("/api/tours/tour/:tourId", authMiddleware, keyPointHandler.GetTourKeyPoints)

	// Key points reorder (posebna grupa)
	r.POST("/api/keypoints/reorder/:tourId", authMiddleware, keyPointHandler.ReorderKeyPoints)

	// Tour durations management
	r.POST("/api/tours/durations", authMiddleware, durationHandler.CreateTourDuration)
	r.DELETE("/api/tours/durations/:id", authMiddleware, durationHandler.DeleteTourDuration)

	// Simulator rute (sa autentifikacijom)
	simulator := r.Group("/api/simulator")
	simulator.Use(authMiddleware)
	{
		simulator.POST("/location", simulatorHandler.UpdateLocation)
		simulator.GET("/location", simulatorHandler.GetCurrentLocation)
		simulator.POST("/execution", simulatorHandler.StartTourExecution)
		simulator.POST("/execution/:id/complete", simulatorHandler.CompleteTourExecution)
	}

	log.Printf("Tour service starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
