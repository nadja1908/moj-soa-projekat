package handler

import (
	"net/http"

	"api-gateway/internal/rpc"

	"github.com/gin-gonic/gin"
)

// RPCHandler handles RPC-based requests
type RPCHandler struct {
	rpcClient *rpc.RPCClient
}

func NewRPCHandler(authAddr, tourAddr string) (*RPCHandler, error) {
	client, err := rpc.NewRPCClient(authAddr, tourAddr)
	if err != nil {
		return nil, err
	}

	return &RPCHandler{rpcClient: client}, nil
}

// Close closes RPC connections
func (h *RPCHandler) Close() {
	if h.rpcClient != nil {
		h.rpcClient.Close()
	}
}

// RPC Login handler
func (h *RPCHandler) RPCLogin(c *gin.Context) {
	var req rpc.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.rpcClient.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RPC Register handler
func (h *RPCHandler) RPCRegister(c *gin.Context) {
	var req rpc.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.rpcClient.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Registration failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// RPC CreateTour handler
func (h *RPCHandler) RPCCreateTour(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req rpc.CreateTourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	args := &rpc.CreateTourArgs{
		UserID:  userID,
		Request: &req,
	}

	resp, err := h.rpcClient.CreateTour(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tour", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// RPC GetTours handler
func (h *RPCHandler) RPCGetTours(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	resp, err := h.rpcClient.GetTours(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tours", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Helper function to extract user ID from context
func getUserIDFromContext(c *gin.Context) int64 {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}

	// Try int64 first
	if id, ok := userID.(int64); ok {
		return id
	}

	// Try int and convert to int64
	if id, ok := userID.(int); ok {
		return int64(id)
	}

	return 0
}
