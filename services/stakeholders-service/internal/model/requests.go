package model

// RegisterRequest je model za parsiranje JSON-a iz zahteva za registraciju
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
}

// LoginRequest je model za parsiranje JSON-a iz zahteva za login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse je model za slanje JWT tokena kao odgovor
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UpdateProfileRequest je model za ažuriranje profila
type UpdateProfileRequest struct {
	FirstName       *string `json:"firstName"`
	LastName        *string `json:"lastName"`
	ProfileImageURL *string `json:"profileImageUrl"`
	Biography       *string `json:"biography"`
	Motto           *string `json:"motto"`
}

type ProfileResponse struct {
    ID              int64  `json:"id"`
    UserID          int64  `json:"userId"`
    FirstName       string `json:"firstName"`
    LastName        string `json:"lastName"`
    ProfileImageURL string `json:"profileImageUrl"`
    Biography       string `json:"biography"`
    Motto           string `json:"motto"`
}