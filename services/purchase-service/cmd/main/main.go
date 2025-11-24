package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"purchase-service/internal/client"
	"purchase-service/internal/handler"
	"purchase-service/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting Purchase Service...")

	port := getEnv("PORT", "8005")
	dbURI := getEnv("DB_URI", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "27017")
	dbUser := getEnv("DB_USER", "")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "purchase_db")

	if dbURI == "" {
		if dbUser != "" && dbPass != "" {
			dbURI = "mongodb://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort
		} else {
			dbURI = "mongodb://" + dbHost + ":" + dbPort
		}
	}

	mongoStore, err := store.NewStore(dbURI, dbName)
	if err != nil {
		log.Fatal("Failed to initialize Mongo store:", err)
	}
	defer mongoStore.Close(context.Background())

	tourClient := client.NewTourClient()

	purchaseHandler := handler.NewPurchaseHandler(mongoStore, tourClient)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	r.GET("/health", purchaseHandler.Health)

	internal := r.Group("/internal")
	{
		internal.GET("/purchase/:touristID", purchaseHandler.InternalGetCart)
	}

	auth := r.Group("/")
	auth.Use(handler.AuthMiddleware())
	{
		auth.POST("/purchase/add", purchaseHandler.AddToCart)
		auth.GET("/purchase", purchaseHandler.GetShoppingCart)
		auth.DELETE("/purchase/:tourId", purchaseHandler.RemoveItem)
		auth.POST("/purchase/checkout", purchaseHandler.Checkout)
	}

	log.Printf("Purchase Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
