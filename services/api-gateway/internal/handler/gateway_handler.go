package handler

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GatewayHandler struct {
	authServiceURL         string
	stakeholdersServiceURL string
	blogServiceURL         string
}

func NewGatewayHandler(authServiceURL, stakeholdersServiceURL, blogServiceURL string) *GatewayHandler {
	return &GatewayHandler{
		authServiceURL:         authServiceURL,
		stakeholdersServiceURL: stakeholdersServiceURL,
		blogServiceURL:         blogServiceURL,
	}
}

func (h *GatewayHandler) ProxyToAuth(c *gin.Context) {
	// Remove /api/auth prefix and proxy to auth service
	path := strings.TrimPrefix(c.Request.URL.Path, "/api/auth")
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

func (h *GatewayHandler) proxyRequest(c *gin.Context, targetURL string) {
	// Read request body
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
	}

	// Create new request
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
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

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to proxy request: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {
			c.Header(name, value)
		}
	}

	// Copy response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Set status and write response
	c.Status(resp.StatusCode)
	c.Writer.Write(respBody)
}