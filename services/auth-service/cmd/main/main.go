package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"auth-service/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := getEnv("PORT", "8003")
	rpcPort := getEnv("RPC_PORT", "9003")
	stakeholdersService := getEnv("STAKEHOLDERS_SERVICE_URL", "http://stakeholders-service:8001")

	// Inicijalizacija handler-a
	authHandler := handler.NewAuthHandler(stakeholdersService)
	authRPCHandler := handler.NewAuthRPCHandler(authHandler)

	// Registracija RPC servera
	rpc.Register(authRPCHandler)
	rpc.HandleHTTP()

	// Pokreni RPC server u goroutine
	go func() {
		log.Printf("Starting RPC server on port %s", rpcPort)
		listener, err := net.Listen("tcp", ":"+rpcPort)
		if err != nil {
			log.Fatalf("Failed to start RPC server: %v", err)
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("RPC connection error: %v", err)
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	// HTTP server setup
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

	log.Printf("Auth service HTTP server starting on port %s", port)
	log.Printf("Auth service RPC server starting on port %s", rpcPort)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
