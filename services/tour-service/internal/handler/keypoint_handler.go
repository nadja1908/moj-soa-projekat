package handler

import (
	"net/http"
	"strconv"

	"tour-service/internal/model"
	"tour-service/internal/store"

	"github.com/gin-gonic/gin"
)

type KeyPointHandler struct {
	store *store.Store
}

func NewKeyPointHandler(store *store.Store) *KeyPointHandler {
	return &KeyPointHandler{store: store}
}

// CreateKeyPoint kreira novu ključnu tačku
func (h *KeyPointHandler) CreateKeyPoint(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req model.CreateKeyPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Proverava da li je korisnik vlasnik ture
	tour, err := h.store.GetTourByID(req.TourID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	if tour.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only add key points to your own tours"})
		return
	}

	// Proverava da li je tura u draft stanju
	if tour.Status != model.StatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only add key points to draft tours"})
		return
	}

	keyPoint, err := h.store.CreateKeyPoint(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create key point"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"keyPoint": keyPoint,
	})
}

// GetKeyPoint vraća ključnu tačku po ID-u
func (h *KeyPointHandler) GetKeyPoint(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key point ID"})
		return
	}

	keyPoint, err := h.store.GetKeyPointByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key point not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"keyPoint": keyPoint,
	})
}

// GetTourKeyPoints vraća sve ključne tačke za turu
func (h *KeyPointHandler) GetTourKeyPoints(c *gin.Context) {
	tourID, err := strconv.ParseInt(c.Param("tourId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	keyPoints, err := h.store.GetKeyPoints(tourID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get key points"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"keyPoints": keyPoints,
	})
}

// UpdateKeyPoint ažurira ključnu tačku
func (h *KeyPointHandler) UpdateKeyPoint(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key point ID"})
		return
	}

	keyPoint, err := h.store.GetKeyPointByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key point not found"})
		return
	}

	// Proverava da li je korisnik vlasnik ture
	tour, err := h.store.GetTourByID(keyPoint.TourID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	if tour.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update key points of your own tours"})
		return
	}

	// Proverava da li je tura u draft stanju
	if tour.Status != model.StatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only update key points of draft tours"})
		return
	}

	var req model.UpdateKeyPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedKeyPoint, err := h.store.UpdateKeyPoint(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update key point"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"keyPoint": updatedKeyPoint,
	})
}

// DeleteKeyPoint briše ključnu tačku
func (h *KeyPointHandler) DeleteKeyPoint(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key point ID"})
		return
	}

	keyPoint, err := h.store.GetKeyPointByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key point not found"})
		return
	}

	// Proverava da li je korisnik vlasnik ture
	tour, err := h.store.GetTourByID(keyPoint.TourID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	if tour.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete key points of your own tours"})
		return
	}

	// Proverava da li je tura u draft stanju
	if tour.Status != model.StatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only delete key points of draft tours"})
		return
	}

	err = h.store.DeleteKeyPoint(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Key point deleted successfully",
	})
}

// ReorderKeyPoints menja redosled ključnih tačaka
func (h *KeyPointHandler) ReorderKeyPoints(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	tourID, err := strconv.ParseInt(c.Param("tourId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	// Proverava da li je korisnik vlasnik ture
	tour, err := h.store.GetTourByID(tourID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	if tour.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only reorder key points of your own tours"})
		return
	}

	// Proverava da li je tura u draft stanju
	if tour.Status != model.StatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only reorder key points of draft tours"})
		return
	}

	var req struct {
		KeyPointIDs []int64 `json:"keyPointIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.store.ReorderKeyPoints(tourID, req.KeyPointIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder key points"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Key points reordered successfully",
	})
}
