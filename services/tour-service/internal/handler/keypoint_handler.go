package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

// calculateTourPrice računa cenu ture na osnovu broja ključnih tačaka
// Osnovna cena: 10€, +5€ za svaku ključnu tačku
func (h *KeyPointHandler) calculateTourPrice(tourID int64) (float64, error) {
	keyPoints, err := h.store.GetKeyPoints(tourID)
	if err != nil {
		return 0, err
	}

	basePrice := 10.0       // Osnovna cena
	pricePerKeyPoint := 5.0 // Cena po ključnoj tački

	totalPrice := basePrice + (float64(len(keyPoints)) * pricePerKeyPoint)
	return totalPrice, nil
}

// updateTourPrice ažurira cenu ture
func (h *KeyPointHandler) updateTourPrice(tourID int64) error {
	price, err := h.calculateTourPrice(tourID)
	if err != nil {
		return err
	}

	// Ažuriraj cenu ture
	updateReq := &model.UpdateTourRequest{
		Price: &price,
	}

	_, err = h.store.UpdateTour(tourID, updateReq)
	return err
}

// CreateKeyPoint kreira novu ključnu tačku
func (h *KeyPointHandler) CreateKeyPoint(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse multipart form
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
		return
	}

	// Extract form fields
	tourID, err := strconv.ParseInt(c.PostForm("tourId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	description := c.PostForm("description")

	latitude, err := strconv.ParseFloat(c.PostForm("latitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	longitude, err := strconv.ParseFloat(c.PostForm("longitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	// Handle image upload
	var imageURL string
	file, header, err := c.Request.FormFile("image")
	if err == nil && file != nil {
		defer file.Close()

		// Create uploads directory if it doesn't exist
		uploadsDir := "./uploads/keypoints"
		os.MkdirAll(uploadsDir, 0755)

		// Generate unique filename
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), name, ext)
		imagePath := filepath.Join(uploadsDir, filename)

		// Save file
		dst, err := os.Create(imagePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		// Set image URL to be served by static file server
		imageURL = fmt.Sprintf("/uploads/keypoints/%s", filename)
	}

	// Create request object
	req := model.CreateKeyPointRequest{
		TourID:      tourID,
		Name:        name,
		Description: description,
		Latitude:    latitude,
		Longitude:   longitude,
		ImageURL:    imageURL,
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

	// Ažuriraj cenu ture na osnovu novih ključnih tačaka
	err = h.updateTourPrice(req.TourID)
	if err != nil {
		log.Printf("Warning: Failed to update tour price: %v", err)
		// Ne prekidamo operaciju zbog greške u ažuriranju cene
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
	log.Printf("DEBUG: GetTourKeyPoints called for tourId: %s", c.Param("tourId"))

	tourID, err := strconv.ParseInt(c.Param("tourId"), 10, 64)
	if err != nil {
		log.Printf("DEBUG: Invalid tour ID error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	log.Printf("DEBUG: Parsed tour ID: %d", tourID)

	keyPoints, err := h.store.GetKeyPoints(tourID)
	if err != nil {
		log.Printf("DEBUG: Store GetKeyPoints error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get key points"})
		return
	}

	log.Printf("DEBUG: Retrieved %d key points", len(keyPoints))

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"keyPoints": keyPoints,
	})
}

// UpdateKeyPoint ažurira ključnu tačku
func (h *KeyPointHandler) UpdateKeyPoint(c *gin.Context) {
	log.Printf("DEBUG: UpdateKeyPoint called for keypoint ID: %s", c.Param("id"))

	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	log.Printf("DEBUG: User ID: %d", userID)

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
		log.Printf("ERROR: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("DEBUG: UpdateKeyPoint request: %+v", req)

	updatedKeyPoint, err := h.store.UpdateKeyPoint(id, &req)
	if err != nil {
		log.Printf("ERROR: Failed to update key point in store: %v", err)
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

	// Ažuriraj cenu ture na osnovu preostalih ključnih tačaka
	err = h.updateTourPrice(keyPoint.TourID)
	if err != nil {
		log.Printf("Warning: Failed to update tour price after deletion: %v", err)
		// Ne prekidamo operaciju zbog greške u ažuriranju cene
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
