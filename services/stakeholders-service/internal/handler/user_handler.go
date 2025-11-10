package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"stakeholders-service/internal/model"
	"stakeholders-service/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	store *store.Store
}

func NewUserHandler(store *store.Store) *UserHandler {
	return &UserHandler{store: store}
}

// Register je handler za registrovanje novih korisnika
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Proveri da li korisničko ime već postoji
	existingUser, err := h.store.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking username"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
		return
	}

	// 2. Proveri da li email već postoji
	existingUser, err = h.store.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking email"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already registered"})
		return
	}

	// 3. Validacija uloge
	switch req.Role {
	case "guide", "tourist":
		// OK
	case "administrator":
		c.JSON(http.StatusForbidden, gin.H{"error": "Administrator role cannot be assigned through registration"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified. Allowed roles are 'guide' or 'tourist'"})
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
	}

	if err := h.store.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// Login je handler za prijavljivanje korisnika
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pronađi korisnika po korisničkom imenu
	user, err := h.store.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Proveri da li je nalog aktivan
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is blocked"})
		return
	}

	// Proveri lozinku
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Kreiraj odgovor bez tokena (Auth servis će kreirati token)
	c.JSON(http.StatusOK, model.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	})
}

// GetAllUsers vraća sve korisnike (samo za administratore)
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "administrator" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Administrator access required"})
		return
	}

	users, err := h.store.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// BlockUser blokira korisnički nalog (samo za administratore)
func (h *UserHandler) BlockUser(c *gin.Context) {
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "administrator" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Administrator access required"})
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.store.BlockUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to block user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User blocked successfully"})
}

// GetUserByID vraća korisnika po ID-u (interno za Auth servis)
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Health check endpoint
func (h *UserHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "stakeholders-service"})
}