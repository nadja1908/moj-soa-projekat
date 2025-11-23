package rpc

import (
	"auth-service/internal/model"
)

// AuthService RPC interface
type AuthService struct {
	handler AuthHandlerInterface
}

type AuthHandlerInterface interface {
	ProcessLogin(req *model.LoginRequest) (*model.LoginResponse, error)
	ProcessRegister(req *model.RegisterRequest) (*model.RegisterResponse, error)
}

// NewAuthService creates a new RPC auth service
func NewAuthService(handler AuthHandlerInterface) *AuthService {
	return &AuthService{handler: handler}
}

// Login RPC method
func (s *AuthService) Login(req *model.LoginRequest, resp *model.LoginResponse) error {
	result, err := s.handler.ProcessLogin(req)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

// Register RPC method
func (s *AuthService) Register(req *model.RegisterRequest, resp *model.RegisterResponse) error {
	result, err := s.handler.ProcessRegister(req)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}
