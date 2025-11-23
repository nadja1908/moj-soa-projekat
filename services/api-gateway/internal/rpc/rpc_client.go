package rpc

import (
	"fmt"
	"net/rpc"
)

// RPCClient wraps RPC connections to various services
type RPCClient struct {
	authConn *rpc.Client
	tourConn *rpc.Client
}

// NewRPCClient creates a new RPC client with connections to services
func NewRPCClient(authAddr, tourAddr string) (*RPCClient, error) {
	client := &RPCClient{}

	// Connect to auth service
	authConn, err := rpc.DialHTTP("tcp", authAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %v", err)
	}
	client.authConn = authConn

	// Connect to tour service
	tourConn, err := rpc.DialHTTP("tcp", tourAddr)
	if err != nil {
		authConn.Close()
		return nil, fmt.Errorf("failed to connect to tour service: %v", err)
	}
	client.tourConn = tourConn

	return client, nil
}

// Close closes all RPC connections
func (c *RPCClient) Close() {
	if c.authConn != nil {
		c.authConn.Close()
	}
	if c.tourConn != nil {
		c.tourConn.Close()
	}
}

// Auth service methods
func (c *RPCClient) Login(req *LoginRequest) (*LoginResponse, error) {
	var resp LoginResponse
	err := c.authConn.Call("AuthService.Login", req, &resp)
	return &resp, err
}

func (c *RPCClient) Register(req *RegisterRequest) (*RegisterResponse, error) {
	var resp RegisterResponse
	err := c.authConn.Call("AuthService.Register", req, &resp)
	return &resp, err
}

// Tour service methods
func (c *RPCClient) CreateTour(args *CreateTourArgs) (*Tour, error) {
	var resp Tour
	err := c.tourConn.Call("TourService.CreateTour", args, &resp)
	return &resp, err
}

func (c *RPCClient) GetTours(userID int64) (*GetToursResponse, error) {
	var resp GetToursResponse
	err := c.tourConn.Call("TourService.GetTours", &userID, &resp)
	return &resp, err
}

// Request/Response types for RPC calls
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refreshToken"`
	User         UserInfo `json:"user"`
}

type RegisterResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	User    UserInfo `json:"user"`
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type CreateTourArgs struct {
	UserID  int64              `json:"userId"`
	Request *CreateTourRequest `json:"request"`
}

type CreateTourRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Difficulty  string   `json:"difficulty"`
	Tags        []string `json:"tags"`
	DistanceKm  float64  `json:"distanceKm"`
}

type Tour struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	AuthorID    int64   `json:"authorId"`
	Difficulty  string  `json:"difficulty"`
	Tags        string  `json:"tags"`
	Status      string  `json:"status"`
	DistanceKm  float64 `json:"distanceKm"`
	Price       float64 `json:"price"`
}

type GetToursResponse struct {
	Tours []*Tour `json:"tours"`
}
