package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"auth-service/internal/model"
)

type AuthRPCHandler struct {
	authHandler *AuthHandler
}

func NewAuthRPCHandler(authHandler *AuthHandler) *AuthRPCHandler {
	return &AuthRPCHandler{
		authHandler: authHandler,
	}
}

func (h *AuthRPCHandler) ProcessLogin(req *model.LoginRPCRequest, resp *model.AuthRPCResponse) error {
	stakeholderURL := "http://stakeholders-service:8001/internal/login"

	loginReq := map[string]string{
		"username": req.Username,
		"password": req.Password,
	}

	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		resp.Error = fmt.Sprintf("Failed to marshal login request: %v", err)
		return nil
	}

	httpResp, err := http.Post(stakeholderURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		resp.Error = fmt.Sprintf("Failed to connect to stakeholders service: %v", err)
		return nil
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Error = fmt.Sprintf("Failed to read response: %v", err)
		return nil
	}

	if httpResp.StatusCode != http.StatusOK {
		resp.Error = fmt.Sprintf("Authentication failed: %s", string(body))
		return nil
	}

	var loginResponse struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		IsActive bool   `json:"isActive"`
	}

	if err := json.Unmarshal(body, &loginResponse); err != nil {
		resp.Error = fmt.Sprintf("Failed to parse login response: %v", err)
		return nil
	}

	user := model.User{
		ID:       loginResponse.ID,
		Username: loginResponse.Username,
		Email:    loginResponse.Email,
		Role:     loginResponse.Role,
		IsActive: loginResponse.IsActive,
	}

	accessToken, refreshToken, expiresAt, err := h.generateTokens(user)
	if err != nil {
		resp.Error = fmt.Sprintf("Failed to generate tokens: %v", err)
		return nil
	}

	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken
	resp.ExpiresAt = expiresAt
	resp.User = user

	log.Printf("RPC Login successful for user: %s", user.Username)
	return nil
}

func (h *AuthRPCHandler) ProcessRegister(req *model.RegisterRPCRequest, resp *model.AuthRPCResponse) error {
	stakeholderURL := "http://stakeholders-service:8001/internal/register"

	registerReq := map[string]string{
		"username": req.Username,
		"password": req.Password,
		"email":    req.Email,
		"role":     req.Role,
	}

	jsonData, err := json.Marshal(registerReq)
	if err != nil {
		resp.Error = fmt.Sprintf("Failed to marshal register request: %v", err)
		return nil
	}

	httpResp, err := http.Post(stakeholderURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		resp.Error = fmt.Sprintf("Failed to connect to stakeholders service: %v", err)
		return nil
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Error = fmt.Sprintf("Failed to read response: %v", err)
		return nil
	}

	if httpResp.StatusCode != http.StatusOK && httpResp.StatusCode != http.StatusCreated {
		resp.Error = fmt.Sprintf("Registration failed: %s", string(body))
		return nil
	}

	var registerResponse struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		IsActive bool   `json:"isActive"`
	}

	if err := json.Unmarshal(body, &registerResponse); err != nil {
		resp.Error = fmt.Sprintf("Failed to parse register response: %v", err)
		return nil
	}

	user := model.User{
		ID:       registerResponse.ID,
		Username: registerResponse.Username,
		Email:    registerResponse.Email,
		Role:     registerResponse.Role,
		IsActive: registerResponse.IsActive,
	}

	accessToken, refreshToken, expiresAt, err := h.generateTokens(user)
	if err != nil {
		resp.Error = fmt.Sprintf("Failed to generate tokens: %v", err)
		return nil
	}

	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken
	resp.ExpiresAt = expiresAt
	resp.User = user

	log.Printf("RPC Register successful for user: %s", user.Username)
	return nil
}

func (h *AuthRPCHandler) generateTokens(user model.User) (string, string, time.Time, error) {
	// Koristimo metode iz postojećeg auth handler-a
	return h.authHandler.GenerateJWT(user.ID, user.Username, user.Role)
}
