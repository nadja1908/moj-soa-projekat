package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GatewayHandler struct {
	authServiceURL         string
	stakeholdersServiceURL string
	blogServiceURL         string
	tourServiceURL         string
}

func NewGatewayHandler(authServiceURL, stakeholdersServiceURL, blogServiceURL, tourServiceURL string) *GatewayHandler {
	return &GatewayHandler{
		authServiceURL:         authServiceURL,
		stakeholdersServiceURL: stakeholdersServiceURL,
		blogServiceURL:         blogServiceURL,
		tourServiceURL:         tourServiceURL,
	}
}

func (h *GatewayHandler) ProxyToAuth(c *gin.Context) {
	// Remove /api/auth prefix and proxy to auth service
	path := strings.TrimPrefix(c.Request.URL.Path, "/api/auth")
	log.Printf("DEBUG: ProxyToAuth - Original path: %s, Final path: %s", c.Request.URL.Path, path)
	log.Printf("DEBUG: ProxyToAuth - Target URL: %s", h.authServiceURL+path)
	h.proxyRequest(c, h.authServiceURL+path)
}

func (h *GatewayHandler) ProxyToStakeholders(c *gin.Context) {
	// Remove /api/users or /api/admin prefix and proxy to stakeholders service
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/api/users") {
		path = strings.TrimPrefix(path, "/api/users")
		if path == "" {
			path = "/profile"
		}
	} else if strings.HasPrefix(path, "/api/admin/users") {
		path = strings.TrimPrefix(path, "/api/admin")
	}

	h.proxyRequest(c, h.stakeholdersServiceURL+path)
}

func (h *GatewayHandler) ProxyToBlog(c *gin.Context) {
	// Remove /api/blog or /api/admin prefix and proxy to blog service
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/api/blog") {
		path = strings.TrimPrefix(path, "/api/blog")
	} else if strings.HasPrefix(path, "/api/admin/posts") {
		path = strings.TrimPrefix(path, "/api/admin/posts")
		path = "/posts" + path
	}

	h.proxyRequest(c, h.blogServiceURL+path)
}

func (h *GatewayHandler) ProxyToTours(c *gin.Context) {
	fmt.Printf("=== PROXY TO TOURS FUNCTION START ===\n")
	log.Printf("DEBUG: ProxyToTours called!")

	// Remove /api/tours prefix and proxy to tour service
	path := strings.TrimPrefix(c.Request.URL.Path, "/api/tours")
	fmt.Printf("=== RAW PATH: '%s' ===\n", path)

	// If path is empty, keep it empty for root endpoint
	if path == "" {
		fmt.Printf("=== PATH WAS EMPTY, KEEPING EMPTY ===\n")
		// path stays empty - will become "/api/tours" without trailing slash
	}

	fmt.Printf("=== FINAL PATH: '%s' ===\n", path)
	log.Printf("DEBUG: ProxyToTours path after fix: %s", path)

	// Get user info from auth middleware context (if available)
	if userIDInterface, exists := c.Get("userID"); exists {
		if userRoleInterface, exists := c.Get("userRole"); exists {
			// Cast from interface{} to proper types
			if userID, ok := userIDInterface.(int); ok {
				if userRole, ok := userRoleInterface.(string); ok {
					// Add user info as headers for tour service
					c.Request.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
					c.Request.Header.Set("X-User-Role", userRole)
					log.Printf("DEBUG: ProxyToTours set headers - X-User-ID: %d, X-User-Role: %s", userID, userRole)
				} else {
					log.Printf("DEBUG: ProxyToTours - userRole cast failed: %T", userRoleInterface)
				}
			} else {
				log.Printf("DEBUG: ProxyToTours - userID cast failed: %T", userIDInterface)
			}
		} else {
			log.Printf("DEBUG: ProxyToTours - no userRole in context")
		}
	} else {
		log.Printf("DEBUG: ProxyToTours - no userID in context")
	}

	// Tour service expects full path with /api/tours prefix
	targetURL := h.tourServiceURL + "/api/tours" + path
	log.Printf("DEBUG: ProxyToTours proxying to: %s", targetURL)
	h.proxyRequest(c, targetURL)
}

func (h *GatewayHandler) proxyRequest(c *gin.Context, targetURL string) {
	log.Printf("DEBUG: proxyRequest called with URL: %s", targetURL)

	// Read request body
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
	}

	// Create new request
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("ERROR: Failed to create proxy request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	// Copy headers
	for name, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// Copy query parameters
	req.URL.RawQuery = c.Request.URL.RawQuery

	log.Printf("DEBUG: Making request to: %s %s", req.Method, req.URL.String())

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: Failed to proxy request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to proxy request: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	log.Printf("DEBUG: Response status: %d", resp.StatusCode)

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {
			c.Header(name, value)
		}
	}

	// Copy response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	log.Printf("DEBUG: Response body: %s", string(respBody))

	// Set status and write response
	c.Status(resp.StatusCode)
	c.Writer.Write(respBody)
}
