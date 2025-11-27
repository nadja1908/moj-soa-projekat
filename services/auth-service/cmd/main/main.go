package main

import (
	"context"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/internal/handler"
	"auth-service/internal/telemetry"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	serviceName := "auth-service"
	port := getEnv("PORT", "8003")
	rpcPort := getEnv("RPC_PORT", "9003")
	stakeholdersService := getEnv("STAKEHOLDERS_SERVICE_URL", "http://stakeholders-service:8001")
	jaegerEndpoint := getEnv("JAEGER_ENDPOINT", "http://jaeger:4317")
	logstashHost := getEnv("LOGSTASH_HOST", "logstash:5000")

	// Initialize structured logging
	logger := telemetry.InitLogger(serviceName, logstashHost)
	logger.Info("Starting auth-service...")

	// Initialize OpenTelemetry tracing
	tp, err := telemetry.InitTracer(serviceName, jaegerEndpoint)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize tracing")
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := telemetry.Shutdown(ctx, tp); err != nil {
			logger.WithError(err).Error("Failed to shutdown tracing")
		}
	}()

	// Inicijalizacija handler-a
	authHandler := handler.NewAuthHandler(stakeholdersService)
	authRPCHandler := handler.NewAuthRPCHandler(authHandler)

	// Registracija RPC servera
	rpc.Register(authRPCHandler)
	rpc.HandleHTTP()

	// Pokreni RPC server u goroutine
	go func() {
		logger.WithField("rpc_port", rpcPort).Info("Starting RPC server")
		listener, err := net.Listen("tcp", ":"+rpcPort)
		if err != nil {
			logger.WithError(err).Fatal("Failed to start RPC server")
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.WithError(err).Error("RPC connection error")
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	// HTTP server setup
	r := gin.Default()

	// Add OpenTelemetry middleware for automatic tracing
	r.Use(otelgin.Middleware(serviceName))

	// Add custom metrics middleware
	r.Use(telemetry.PrometheusMetrics())

	// CORS konfiguracija
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// Setup metrics endpoint
	telemetry.SetupMetricsEndpoint(r)

	// Health check
	r.GET("/health", authHandler.Health)

	// Auth endpoints
	r.POST("/login", authHandler.Login)
	r.POST("/register", authHandler.Register)
	r.POST("/verify", authHandler.VerifyToken)
	r.POST("/refresh", authHandler.RefreshToken)

	// Token validation endpoint za druge servise
	r.GET("/validate", authHandler.ValidateToken)

	// Setup graceful shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		logger.WithField("port", port).Info("Auth service HTTP server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down auth-service...")

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Failed to shutdown HTTP server")
	}

	logger.Info("Auth-service stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
