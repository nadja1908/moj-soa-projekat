package handler

import (
	"fmt"
	"log"
	"net/rpc"
	"time"

	"github.com/gin-gonic/gin"
)

type RPCHandler struct {
	authRPCClient  *rpc.Client
	tourRPCClient  *rpc.Client
	gatewayHandler *GatewayHandler
}

// RPC modeli (moraju biti identicni sa auth service)
type LoginRPCRequest struct {
	Username string
	Password string
}

type RegisterRPCRequest struct {
	Username string
	Password string
	Email    string
	Role     string
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
}

type AuthRPCResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	User         User
	Error        string
}

func NewRPCHandler(authAddr, tourAddr string, gatewayHandler *GatewayHandler) (*RPCHandler, error) {
	// Konektovanje na auth RPC servis
	authClient, err := rpc.Dial("tcp", authAddr)
	if err != nil {
		log.Printf("Error connecting to auth RPC at %s: %v", authAddr, err)
		return &RPCHandler{gatewayHandler: gatewayHandler}, nil // Vraćamo prazan handler - RPC neće raditi ali neće pucati app
	}

	log.Printf("Connected to Auth RPC at %s", authAddr)

	// Za tour service koristimo gRPC kroz gateway handler
	return &RPCHandler{
		authRPCClient:  authClient,
		tourRPCClient:  nil, // Tour service koristi gRPC
		gatewayHandler: gatewayHandler,
	}, nil
}

func (h *RPCHandler) LoginRPC(c *gin.Context) {
	if h.authRPCClient == nil {
		c.JSON(500, gin.H{"error": "Auth RPC service not available"})
		return
	}

	var req LoginRPCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	var resp AuthRPCResponse
	err := h.authRPCClient.Call("AuthRPCHandler.ProcessLogin", &req, &resp)
	if err != nil {
		log.Printf("RPC Login error: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("RPC call failed: %v", err)})
		return
	}

	if resp.Error != "" {
		c.JSON(401, gin.H{"error": resp.Error})
		return
	}

	c.JSON(200, gin.H{
		"accessToken":  resp.AccessToken,
		"refreshToken": resp.RefreshToken,
		"expiresAt":    resp.ExpiresAt,
		"user":         resp.User,
	})
}

func (h *RPCHandler) RegisterRPC(c *gin.Context) {
	if h.authRPCClient == nil {
		c.JSON(500, gin.H{"error": "Auth RPC service not available"})
		return
	}

	var req RegisterRPCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	var resp AuthRPCResponse
	err := h.authRPCClient.Call("AuthRPCHandler.ProcessRegister", &req, &resp)
	if err != nil {
		log.Printf("RPC Register error: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("RPC call failed: %v", err)})
		return
	}

	if resp.Error != "" {
		c.JSON(400, gin.H{"error": resp.Error})
		return
	}

	c.JSON(201, gin.H{
		"accessToken":  resp.AccessToken,
		"refreshToken": resp.RefreshToken,
		"expiresAt":    resp.ExpiresAt,
		"user":         resp.User,
	})
}

func (h *RPCHandler) Close() error {
	if h.authRPCClient != nil {
		return h.authRPCClient.Close()
	}
	if h.tourRPCClient != nil {
		return h.tourRPCClient.Close()
	}
	return nil
}

// Tour RPC methods using gRPC through gateway handler
func (h *RPCHandler) RPCCreateTour(c *gin.Context) {
	// TODO: Implementirati create tour preko gRPC
	c.JSON(501, gin.H{"error": "Tour RPC create not implemented yet"})
}

func (h *RPCHandler) RPCGetTours(c *gin.Context) {
	// Koristimo postojeći gRPC poziv kroz gateway handler
	if h.gatewayHandler == nil {
		c.JSON(500, gin.H{"error": "Gateway handler not available"})
		return
	}

	// Pozivamo postojeću GetPublishedToursRPC metodu
	h.gatewayHandler.GetPublishedToursRPC(c)
}
