package handler

import (
	"bytes"
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
}

func NewGatewayHandler(authServiceURL, stakeholdersServiceURL, blogServiceURL string) *GatewayHandler {
	return &GatewayHandler{
		authServiceURL:         authServiceURL,
		stakeholdersServiceURL: stakeholdersServiceURL,
		blogServiceURL:         blogServiceURL,
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

	log.Println("[ProxyToAuth] PATH:", path)
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

	// Javne blogg rute
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
// MASTER PROXY LOGIKA
// ////////////////////////
func (h *GatewayHandler) proxyRequest(c *gin.Context, targetURL string) {

	// 1) UČITAJ BODY IZ REQUESTA
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body.Close()

	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(bodyBytes))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create proxy request"})
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

	// 4) POŠALJI DALJE SERVISU
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Proxy error: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	// 5) PREKOPIRAJ HEADER-e NAZAD KLIJENTU
	for k, v := range resp.Header {
		for _, val := range v {
			c.Header(k, val)
		}
	}

	// 6) PREKOPIRAJ BODY
	respBody, _ := io.ReadAll(resp.Body)

	// DEBUG — ŠTA BLOG-SERVICE VRAĆA
	log.Printf("[Proxy RESPONSE] Status=%d Body=%s", resp.StatusCode, string(respBody))
	log.Println("------------------------------------------------")

	c.Status(resp.StatusCode)
	c.Writer.Write(respBody)
}
