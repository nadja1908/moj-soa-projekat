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

// ////////////////////////
// AUTH SERVICE PROXY
// ////////////////////////
func (h *GatewayHandler) ProxyToAuth(c *gin.Context) {
	path := strings.TrimPrefix(c.Request.URL.Path, "/api/auth")
	if path == "" {
		path = "/"
	}

	log.Printf("DEBUG: ProxyToAuth - Original path: %s, Final path: %s", c.Request.URL.Path, path)
	log.Printf("DEBUG: ProxyToAuth - Target URL: %s", h.authServiceURL+path)
	h.proxyRequest(c, h.authServiceURL+path)
}

// ////////////////////////
// STAKEHOLDERS SERVICE PROXY
// ////////////////////////
func (h *GatewayHandler) ProxyToStakeholders(c *gin.Context) {
	original := c.Request.URL.Path
	path := original

	if strings.HasPrefix(original, "/api/users") {
		path = strings.TrimPrefix(original, "/api/users")
		if path == "" {
			path = "/profile"
		}
	}

	if strings.HasPrefix(original, "/api/admin/users") {
		path = strings.TrimPrefix(original, "/api/admin")
	}

	final := h.stakeholdersServiceURL + path
	log.Println("[ProxyToStakeholders] →", final)
	h.proxyRequest(c, final)
}

// ////////////////////////
// BLOG SERVICE PROXY
// ////////////////////////
func (h *GatewayHandler) ProxyToBlog(c *gin.Context) {
	// DEBUG: Šta je STIGLO IZ BROWSER-a
	log.Println("------------------------------------------------")
	log.Println("[ProxyToBlog] ORIGINAL PATH:", c.Request.URL.Path)
	log.Println("[ProxyToBlog] AUTH FROM CLIENT:", c.GetHeader("Authorization"))

	original := c.Request.URL.Path
	path := original

	// Javne blog rute
	if strings.HasPrefix(original, "/api/blog") {
		path = strings.TrimPrefix(original, "/api/blog")
		if path == "" {
			path = "/"
		}
	}

	// Admin blog rute
	if strings.HasPrefix(original, "/api/admin/posts") {
		path = strings.TrimPrefix(original, "/api/admin")
		if !strings.HasPrefix(path, "/posts") {
			path = "/posts" + path
		}
	}

	finalURL := h.blogServiceURL + path
	log.Println("[ProxyToBlog] FINAL URL:", finalURL)

	h.proxyRequest(c, finalURL)
}

// ////////////////////////
// TOUR SERVICE PROXY
// ////////////////////////
func (h *GatewayHandler) ProxyToTours(c *gin.Context) {
	log.Printf("DEBUG: ProxyToTours called!")

	// Remove /api/tours prefix and proxy to tour service
	path := strings.TrimPrefix(c.Request.URL.Path, "/api/tours")

	// If path is empty, keep it empty for root endpoint
	if path == "" {
		// path stays empty - will become "/api/tours" without trailing slash
	}

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
				}
			}
		}
	}

	// Tour service expects full path with /api/tours prefix
	targetURL := h.tourServiceURL + "/api/tours" + path
	log.Printf("DEBUG: ProxyToTours proxying to: %s", targetURL)
	h.proxyRequest(c, targetURL)
}

// ////////////////////////
// MASTER PROXY LOGIKA
// ////////////////////////
func (h *GatewayHandler) proxyRequest(c *gin.Context, targetURL string) {
	log.Printf("DEBUG: proxyRequest called with URL: %s", targetURL)

	// 1) UČITAJ BODY IZ REQUESTA
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body.Close()

	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("ERROR: Failed to create proxy request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	// 2) PRENESI SVE HEADER-e (uključujući AUTH)
	for k, v := range c.Request.Header {
		for _, val := range v {
			req.Header.Add(k, val)
		}
	}

	// DEBUG — DA VIDIMO DA LI JE AUTH ZAISTA PROSLEĐEN
	log.Println("[Proxy] SENDING →", targetURL)
	log.Println("[Proxy] AUTH FORWARDED:", req.Header.Get("Authorization"))

	// 3) Query parametri
	req.URL.RawQuery = c.Request.URL.RawQuery

	log.Printf("DEBUG: Making request to: %s %s", req.Method, req.URL.String())

	// 4) POŠALJI DALJE SERVISU
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("ERROR: Failed to proxy request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proxy error: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	log.Printf("DEBUG: Response status: %d", resp.StatusCode)

	// 5) PREKOPIRAJ HEADER-e NAZAD KLIJENTU
	for k, v := range resp.Header {
		for _, val := range v {
			c.Header(k, val)
		}
	}

	// 6) PREKOPIRAJ BODY
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// DEBUG — ŠTA SERVICE VRAĆA
	log.Printf("[Proxy RESPONSE] Status=%d Body=%s", resp.StatusCode, string(respBody))
	log.Println("------------------------------------------------")

	// Set status and write response
	c.Status(resp.StatusCode)
	c.Writer.Write(respBody)
}

// ////////////////////////
// KEYPOINTS PROXY
// ////////////////////////
func (h *GatewayHandler) ProxyToKeyPoints(c *gin.Context) {
	log.Printf("DEBUG: ProxyToKeyPoints called!")

	// Remove /api/keypoints prefix and map to tour service keypoints routes
	path := strings.TrimPrefix(c.Request.URL.Path, "/api/keypoints")

	log.Printf("DEBUG: ProxyToKeyPoints path after trim: %s", path)

	// Get user info from auth middleware context (if available)
	if userIDInterface, exists := c.Get("userID"); exists {
		if userRoleInterface, exists := c.Get("userRole"); exists {
			// Cast from interface{} to proper types
			if userID, ok := userIDInterface.(int); ok {
				if userRole, ok := userRoleInterface.(string); ok {
					// Add user info as headers for tour service
					c.Request.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
					c.Request.Header.Set("X-User-Role", userRole)
					log.Printf("DEBUG: ProxyToKeyPoints set headers - X-User-ID: %d, X-User-Role: %s", userID, userRole)
				}
			}
		}
	}

	// Map keypoints routes to tour service routes
	var targetPath string
	if path == "" {
		// POST /api/keypoints -> /api/tours/keypoints
		targetPath = "/api/tours/keypoints"
	} else if strings.HasPrefix(path, "/tour/") {
		// GET /api/keypoints/tour/123 -> /api/keypoints/tour/123 (već postoji u tour service)
		targetPath = "/api/keypoints" + path
	} else if strings.HasPrefix(path, "/reorder/") {
		// POST /api/keypoints/reorder/123 -> /api/keypoints/reorder/123 (već postoji u tour service)
		targetPath = "/api/keypoints" + path
	} else {
		// GET/PUT/DELETE /api/keypoints/123 -> /api/tours/keypoints/123
		targetPath = "/api/tours/keypoints" + path
	}

	targetURL := h.tourServiceURL + targetPath
	log.Printf("DEBUG: ProxyToKeyPoints proxying to: %s", targetURL)
	h.proxyRequest(c, targetURL)
}
