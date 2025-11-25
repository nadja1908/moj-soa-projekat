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

	// Inicijalizacija STORE-a (glavni store sa konekcijom)
	mainStore := store.NewStore(dbUser, dbPass, dbHost, dbName)
	defer mainStore.Close()

	// Inicijalizacija pojedinačnih store-ova
	tourStore := mainStore
	keyPointStore := mainStore
	durationStore := mainStore
	simulatorStore := mainStore

	// *** NOVO — ReviewStore (ovde koristimo GetDB() metod) ***
	reviewStore := store.NewReviewStore(mainStore.GetDB())

	// Handleri
	tourHandler := handler.NewTourHandler(tourStore)
	keyPointHandler := handler.NewKeyPointHandler(keyPointStore)
	durationHandler := handler.NewTourDurationHandler(durationStore)
	simulatorHandler := handler.NewTourSimulatorHandler(simulatorStore)

	// *** NOVO — ReviewHandler ***
	reviewHandler := handler.NewReviewHandler(reviewStore)

	// Gin setup
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-User-ID", "X-User-Role"},
		AllowCredentials: true,
	}))

	// JWT middleware (API Gateway prosleđuje X-User-ID i X-User-Role)
	authMiddleware := func(c *gin.Context) {
		userIDHeader := c.GetHeader("X-User-ID")

		if userIDHeader == "" {
			c.JSON(401, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userID, err := strconv.ParseInt(userIDHeader, 10, 64)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		c.Set("userId", userID)
		c.Next()
	}

	// Static za slike
	r.Static("/uploads", "./uploads")

	// Health
	r.GET("/health", tourHandler.Health)

	// Public tours
	r.GET("/api/tours/published", tourHandler.GetPublishedTours)
	r.GET("/api/tours/public/:id", tourHandler.GetTourForTourist)

	// Public keypoints + durations
	r.GET("/api/keypoints/tour/:tourId", keyPointHandler.GetTourKeyPoints)
	r.GET("/api/durations/tour/:tourId", durationHandler.GetTourDurations)

	// Protected tours
	r.POST("/api/tours", authMiddleware, tourHandler.CreateTour)
	r.GET("/api/tours/my", authMiddleware, tourHandler.GetMyTours)
	r.PUT("/api/tours/:id", authMiddleware, tourHandler.UpdateTour)
	r.POST("/api/tours/:id/publish", authMiddleware, tourHandler.PublishTour)
	r.POST("/api/tours/:id/archive", authMiddleware, tourHandler.ArchiveTour)
	r.POST("/api/tours/:id/reactivate", authMiddleware, tourHandler.ReactivateTour)
	r.DELETE("/api/tours/:id", authMiddleware, tourHandler.DeleteTour)

	// Keypoints
	r.POST("/api/tours/keypoints", authMiddleware, keyPointHandler.CreateKeyPoint)
	r.GET("/api/tours/keypoints/:id", authMiddleware, keyPointHandler.GetKeyPoint)
	r.PUT("/api/tours/keypoints/:id", authMiddleware, keyPointHandler.UpdateKeyPoint)
	r.DELETE("/api/tours/keypoints/:id", authMiddleware, keyPointHandler.DeleteKeyPoint)
	r.GET("/api/tours/tour/:tourId", authMiddleware, keyPointHandler.GetTourKeyPoints)

	// Reorder
	r.POST("/api/keypoints/reorder/:tourId", authMiddleware, keyPointHandler.ReorderKeyPoints)

	// Durations
	r.POST("/api/tours/durations", authMiddleware, durationHandler.CreateTourDuration)
	r.DELETE("/api/tours/durations/:id", authMiddleware, durationHandler.DeleteTourDuration)

	// Simulator
	simulator := r.Group("/api/simulator")
	simulator.Use(authMiddleware)
	{
		simulator.POST("/location", simulatorHandler.UpdateLocation)
		simulator.GET("/location", simulatorHandler.GetCurrentLocation)
		simulator.POST("/execution", simulatorHandler.StartTourExecution)
		simulator.POST("/execution/:id/complete", simulatorHandler.CompleteTourExecution)
	}

	// *** REVIEWS ***
	// Public:
	r.GET("/api/reviews/tour/:tourId", reviewHandler.GetReviewsByTourID)

	// Protected:
	r.POST("/api/reviews", authMiddleware, reviewHandler.CreateReview)

	log.Printf("Tour service starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
