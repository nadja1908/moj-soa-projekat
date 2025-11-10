package model

import "time"

// LoginRequest je model za parsiranje JSON-a iz zahteva za login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest je model za parsiranje JSON-a iz zahteva za registraciju
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
}

// AuthResponse je model za slanje JWT tokena kao odgovor
type AuthResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	User         User      `json:"user"`
}

// TokenValidationRequest je model za validaciju tokena
type TokenValidationRequest struct {
	Token string `json:"token" binding:"required"`
}

// TokenValidationResponse je odgovor na validaciju tokena
type TokenValidationResponse struct {
	Valid  bool   `json:"valid"`
	UserID int64  `json:"userId,omitempty"`
	Role   string `json:"role,omitempty"`
	Error  string `json:"error,omitempty"`
}

// RefreshTokenRequest je model za refresh token zahtev
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// User predstavlja korisničke podatke
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
}