package telemetry

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PrometheusMetrics middleware for collecting HTTP metrics
func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log metrics
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// For now, we'll log metrics in a format that can be parsed by Prometheus
		// In a real implementation, you would use prometheus/client_golang
		logrus.WithFields(logrus.Fields{
			"method":      method,
			"path":        path,
			"status_code": status,
			"duration_ms": duration.Milliseconds(),
			"user_agent":  c.Request.UserAgent(),
			"remote_addr": c.ClientIP(),
			"metric_type": "http_request",
		}).Info("HTTP request processed")
	}
}

// SetupMetricsEndpoint adds /metrics endpoint for Prometheus scraping
func SetupMetricsEndpoint(router *gin.Engine) {
	router.GET("/metrics", gin.HandlerFunc(func(c *gin.Context) {
		// Basic metrics endpoint - in production you'd use prometheus/client_golang
		c.String(http.StatusOK, `# HELP auth_service_info Service information
# TYPE auth_service_info gauge
auth_service_info{version="1.0.0",service="auth-service"} 1
# HELP auth_service_up Service availability
# TYPE auth_service_up gauge
auth_service_up 1
`)
	}))
}
