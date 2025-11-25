package handler

import (
	"bytes"
	"encoding/json"
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
	followerServiceURL     string
}

func NewGatewayHandler(authServiceURL, stakeholdersServiceURL, blogServiceURL, tourServiceURL, followerServiceURL string, purchaseURL string) *GatewayHandler {
	purchaseServiceURL     string
}

	return &GatewayHandler{
		authServiceURL:         authServiceURL,
		stakeholdersServiceURL: stakeholdersServiceURL,
		blogServiceURL:         blogServiceURL,
		tourServiceURL:         tourServiceURL,
		followerServiceURL:     followerServiceURL,
		purchaseServiceURL:     purchaseURL,
	}
}

func (h *GatewayHandler) ProxyToAuth(c *gin.Context) {
	path := strings.TrimPrefix(c.Request.URL.Path, "/api/auth")
	if path == "" {
		path = "/"
	}

	log.Printf("DEBUG: ProxyToAuth - Original path: %s, Final path: %s", c.Request.URL.Path, path)
	log.Printf("DEBUG: ProxyToAuth - Target URL: %s", h.authServiceURL+path)
	h.proxyRequest(c, h.authServiceURL+path)
}

func (h *GatewayHandler) ProxyToPurchase(c *gin.Context) {
	original := c.Request.URL.Path
	path := original

	if path == "" {
		path = "/"
	}
	if strings.HasPrefix(original, "/api/purchase") {
		path = strings.TrimPrefix(original, "/api/purchase")
		if path == "" {
			path = "/"
		}
	}

	finalURL := h.purchaseServiceURL + "/purchase" + path
	log.Println("[ProxyToPurchase] FINAL URL →", finalURL)
	h.proxyRequest(c, finalURL)
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
// FOLLOWER SERVICE PROXY
// ////////////////////////
func (h *GatewayHandler) ProxyToFollower(c *gin.Context) {
	log.Println("------------------------------------------------")
	log.Println("[ProxyToFollower] ORIGINAL PATH:", c.Request.URL.Path)
	log.Println("[ProxyToFollower] METHOD:", c.Request.Method)
	log.Println("[ProxyToFollower] AUTH FROM CLIENT:", c.GetHeader("Authorization"))

	original := c.Request.URL.Path
	path := strings.TrimPrefix(original, "/api/follower")
	if path == "" {
		path = "/"
	}

	// Follower service očekuje /api/follower prefix
	finalURL := h.followerServiceURL + "/api/follower" + path
	log.Println("[ProxyToFollower] FINAL URL:", finalURL)

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

// ////////////////////////
// SAGA ORCHESTRATION FOR RECOMMENDATIONS
// ////////////////////////

// UserRecommendation struktura za preporuke korisnika
type UserRecommendation struct {
	UserId          int    `json:"userId"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	FirstName       string `json:"firstName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	Role            string `json:"role,omitempty"`
	CommonFollowers int    `json:"commonFollowers"`
}

// StakeholderDetails struktura za detalje korisnika iz stakeholders servisa
type StakeholderDetails struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
}

type StakeholderResponse struct {
	User StakeholderDetails `json:"user"`
}

// GetRecommendationsWithDetails implementira SAGA orchestration pattern
// Koraci:
// 1. Poziva follower-service za preporuke iz Neo4j
// 2. Za svaki userId poziva stakeholders-service da dobije dodatne detalje
// 3. Obogaćuje podatke i vraća kompletan rezultat
// 4. Ako stakeholders service ne odgovori, vraća osnovne podatke (fallback/compensation)
func (h *GatewayHandler) GetRecommendationsWithDetails(c *gin.Context) {
	log.Println("========================================")
	log.Println("SAGA ORCHESTRATION: Starting recommendations saga")
	log.Println("========================================")

	// Korak 1: Pozovi follower service za osnovne preporuke
	log.Println("SAGA Step 1: Calling follower-service for recommendations...")
	
	// Kreiraj HTTP zahtev ka follower servisu
	followerURL := h.followerServiceURL + "/api/follower/recommendations"
	req, err := http.NewRequest("GET", followerURL, nil)
	if err != nil {
		log.Printf("SAGA ERROR: Failed to create request to follower service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Prekopiraj Authorization header
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	// Pozovi follower service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("SAGA ERROR: Failed to call follower service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("SAGA ERROR: Follower service returned status %d", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		log.Printf("SAGA ERROR: Response body: %s", string(body))
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to get recommendations from follower service"})
		return
	}

	// Parse osnovne preporuke
	var basicRecommendations []UserRecommendation
	if err := json.NewDecoder(resp.Body).Decode(&basicRecommendations); err != nil {
		log.Printf("SAGA ERROR: Failed to decode recommendations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse recommendations"})
		return
	}

	log.Printf("SAGA Step 1 SUCCESS: Received %d recommendations from follower service", len(basicRecommendations))

	// Korak 2: Za svaki userId, pozovi stakeholders service da obogatiš podatke
	log.Println("SAGA Step 2: Enriching recommendations with stakeholder details...")
	
	enrichedRecommendations := make([]UserRecommendation, 0, len(basicRecommendations))
	successCount := 0
	failureCount := 0

	for i, rec := range basicRecommendations {
		log.Printf("SAGA Step 2.%d: Fetching details for userId=%d (username=%s)", i+1, rec.UserId, rec.Username)
		
		// Pozovi stakeholders service za detalje
		stakeholderURL := fmt.Sprintf("%s/internal/users/%d", h.stakeholdersServiceURL, rec.UserId)
		log.Printf("SAGA DEBUG: Full stakeholder URL: %s", stakeholderURL)
		stakeholderReq, err := http.NewRequest("GET", stakeholderURL, nil)
		if err != nil {
			log.Printf("SAGA COMPENSATION: Failed to create request for userId=%d, using basic data", rec.UserId)
			enrichedRecommendations = append(enrichedRecommendations, rec)
			failureCount++
			continue
		}

		stakeholderResp, err := client.Do(stakeholderReq)
		if err != nil {
			log.Printf("SAGA COMPENSATION: Failed to call stakeholders service for userId=%d, using basic data: %v", rec.UserId, err)
			enrichedRecommendations = append(enrichedRecommendations, rec)
			failureCount++
			continue
		}

		log.Printf("SAGA DEBUG: Stakeholders response status: %d for userId=%d", stakeholderResp.StatusCode, rec.UserId)
		if stakeholderResp.StatusCode != http.StatusOK {
			log.Printf("SAGA COMPENSATION: Stakeholders service returned status %d for userId=%d, using basic data", stakeholderResp.StatusCode, rec.UserId)
			stakeholderResp.Body.Close()
			enrichedRecommendations = append(enrichedRecommendations, rec)
			failureCount++
			continue
		}

		// Parse detalje o korisniku
		var stakeholderResponse StakeholderResponse
		if err := json.NewDecoder(stakeholderResp.Body).Decode(&stakeholderResponse); err != nil {
			log.Printf("SAGA COMPENSATION: Failed to parse stakeholder details for userId=%d, using basic data: %v", rec.UserId, err)
			stakeholderResp.Body.Close()
			enrichedRecommendations = append(enrichedRecommendations, rec)
			failureCount++
			continue
		}
		stakeholderResp.Body.Close()

		// Obogati preporuku sa detaljima
		stakeholderDetails := stakeholderResponse.User
		rec.Email = stakeholderDetails.Email
		rec.FirstName = stakeholderDetails.FirstName
		rec.LastName = stakeholderDetails.LastName
		rec.Role = stakeholderDetails.Role
		
		enrichedRecommendations = append(enrichedRecommendations, rec)
		successCount++
		log.Printf("SAGA Step 2.%d SUCCESS: Enriched userId=%d with email=%s, role=%s", i+1, rec.UserId, rec.Email, rec.Role)
	}

	log.Println("========================================")
	log.Printf("SAGA ORCHESTRATION COMPLETED:")
	log.Printf("  - Total recommendations: %d", len(basicRecommendations))
	log.Printf("  - Successfully enriched: %d", successCount)
	log.Printf("  - Fallback to basic data: %d", failureCount)
	log.Println("========================================")

	// Korak 3: Vrati obogaćene preporuke
	c.JSON(http.StatusOK, enrichedRecommendations)
}

