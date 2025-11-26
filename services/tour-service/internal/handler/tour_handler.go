package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"tour-service/internal/model"
	"tour-service/internal/store"
	tourpb "tour-service/proto"

	"github.com/gin-gonic/gin"
)

type TourHandler struct {
	store *store.Store
}

func NewTourHandler(store *store.Store) *TourHandler {
	return &TourHandler{store: store}
}

// CreateTour kreira novu turu
func (h *TourHandler) CreateTour(c *gin.Context) {
	log.Printf("DEBUG: CreateTour STARTED!")

	userID := getUserIDFromContext(c)
	log.Printf("DEBUG: CreateTour - userID from context: %d", userID)

	if userID == 0 {
		log.Printf("DEBUG: CreateTour - userID is 0, returning 401")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req model.CreateTourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("DEBUG: CreateTour - JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("DEBUG: CreateTour - Request: %+v", req)

	tour, err := h.store.CreateTour(userID, &req)
	if err != nil {
		log.Printf("DEBUG: CreateTour - Store error: %v", err)
		log.Printf("DEBUG: CreateTour - Store error type: %T", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tour"})
		return
	}

	log.Printf("DEBUG: CreateTour - Tour created successfully: %+v", tour)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"tour":    tour,
	})
}

// GetTour vraća turu po ID-u
func (h *TourHandler) GetTour(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	tour, err := h.store.GetTourByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tour":    tour,
	})
}

// GetMyTours vraća sve ture trenutnog korisnika (autor)
func (h *TourHandler) GetMyTours(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	tours, err := h.store.GetToursByAuthor(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tours"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tours":   tours,
	})
}

// GetPublishedTours vraća sve objavljene ture (za turiste)
func (h *TourHandler) GetPublishedTours(c *gin.Context) {
	tours, err := h.store.GetPublishedTours()
	if err != nil {
		log.Printf("Error getting published tours: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get published tours"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tours":   tours,
	})
}

// GetTourForTourist vraća objavljenu turu za turiste (sa samo prvom ključnom tačkom)
func (h *TourHandler) GetTourForTourist(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	log.Printf("DEBUG: GetTourForTourist called for tour ID: %d", id)

	tour, err := h.store.GetTourForTourist(id)
	if err != nil {
		log.Printf("DEBUG: GetTourForTourist error: %v", err)
		if err.Error() == "published tour not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found or not published"})
		} else {
			log.Printf("Error getting tour for tourist: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tour"})
		}
		return
	}

	log.Printf("DEBUG: GetTourForTourist success: %+v", tour)

	tags := []string{}
	if tour.Tour.Tags != "" {
		_ = json.Unmarshal([]byte(tour.Tour.Tags), &tags)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tour": gin.H{
			"id":          tour.Tour.ID,
			"name":        tour.Tour.Name,
			"price":       tour.Tour.Price,
			"status":      string(tour.Tour.Status),
			"description": tour.Tour.Description,
			"distanceKm":  tour.Tour.DistanceKm,
			"difficulty":  string(tour.Tour.Difficulty),
			"tags":        tags,
			"publishedAt": tour.Tour.PublishedAt,
		},
	})
}

// UpdateTour ažurira turu
func (h *TourHandler) UpdateTour(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	// Proverava da li je korisnik vlasnik ture
	tour, err := h.store.GetTourByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	if tour.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own tours"})
		return
	}

	var req model.UpdateTourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTour, err := h.store.UpdateTour(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tour":    updatedTour,
	})
}

// PublishTour objavljuje turu
func (h *TourHandler) PublishTour(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	// Proverava da li je korisnik vlasnik ture
	tour, err := h.store.GetTourByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	if tour.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only publish your own tours"})
		return
	}

	var req model.PublishTourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.store.PublishTour(id, req.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTour, err := h.store.GetTourByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated tour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tour":    updatedTour,
		"message": "Tour published successfully",
	})
}

// ArchiveTour arhivira turu
func (h *TourHandler) ArchiveTour(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	// Proverava da li je korisnik vlasnik ture
	tour, err := h.store.GetTourByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	if tour.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only archive your own tours"})
		return
	}

	err = h.store.ArchiveTour(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tour archived successfully",
	})
}

// ReactivateTour ponovo aktivira arhiviranu turu
func (h *TourHandler) ReactivateTour(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	// Proverava da li je korisnik vlasnik ture
	tour, err := h.store.GetTourByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	if tour.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only reactivate your own tours"})
		return
	}

	err = h.store.ReactivateTour(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tour reactivated successfully",
	})
}

// DeleteTour briše turu (draft i archived ture)
func (h *TourHandler) DeleteTour(c *gin.Context) {
	log.Printf("🗑️ DELETE request received for tour")

	userID := getUserIDFromContext(c)
	log.Printf("🔍 User ID from context: %d", userID)
	if userID == 0 {
		log.Printf("❌ User not authenticated")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("❌ Invalid tour ID: %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	log.Printf("🎯 Attempting to delete tour ID: %d for user: %d", id, userID)

	err = h.store.DeleteTour(id, userID)
	if err != nil {
		log.Printf("❌ Delete failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("✅ Tour %d deleted successfully", id)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tour deleted successfully",
	})
}

// Health check
func (h *TourHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "tour-service",
	})
}

// getUserIDFromContext izvlači user ID iz context-a (postavlja middleware)
func getUserIDFromContext(c *gin.Context) int64 {
	userIDInterface, exists := c.Get("userId")
	if !exists {
		log.Printf("DEBUG: No userId in context")
		return 0
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Printf("DEBUG: userId in context is not int64: %T", userIDInterface)
		return 0
	}

	log.Printf("DEBUG: Got userId from context: %d", userID)
	return userID
}

// RPC endpoints - HTTP interface to internal RPC calls
func (h *TourHandler) GetPublishedToursRPC(c *gin.Context) {
	log.Printf("DEBUG: GetPublishedToursRPC called - HTTP endpoint that calls internal RPC")

	// Pozivamo postojeći TourRPCServer direktno
	rpcServer := NewTourRPCServer(h.store)

	// Kreiraj gRPC request
	resp, err := rpcServer.GetPublishedTours(c.Request.Context(), &tourpb.GetPublishedToursRequest{})
	if err != nil {
		log.Printf("ERROR: RPC GetPublishedTours failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tours via RPC"})
		return
	}

	// Mapiranje na HTTP odgovor
	tours := make([]gin.H, 0, len(resp.Tours))
	for _, t := range resp.Tours {
		tours = append(tours, gin.H{
			"id":          t.Id,
			"name":        t.Name,
			"description": t.Description,
			"difficulty":  t.Difficulty,
			"status":      t.Status,
			"price":       t.Price,
			"distanceKm":  t.DistanceKm,
			"publishedAt": t.PublishedAt,
			"tags":        t.Tags,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tours":   tours,
	})
}

func (h *TourHandler) GetTourForTouristRPC(c *gin.Context) {
	log.Printf("DEBUG: GetTourForTouristRPC called - HTTP endpoint that calls internal RPC")

	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing tour id"})
		return
	}

	var id int64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour id"})
		return
	}

	// Pozivamo postojeći TourRPCServer direktno
	rpcServer := NewTourRPCServer(h.store)

	// Kreiraj gRPC request
	resp, err := rpcServer.GetPublicTour(c.Request.Context(), &tourpb.GetPublicTourRequest{Id: id})
	if err != nil {
		log.Printf("ERROR: RPC GetPublicTour failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tour via RPC"})
		return
	}

	if resp.Tour == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	t := resp.Tour

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tour": gin.H{
			"id":          t.Id,
			"name":        t.Name,
			"price":       t.Price,
			"status":      t.Status,
			"description": t.Description,
			"distanceKm":  t.DistanceKm,
			"difficulty":  t.Difficulty,
			"tags":        t.Tags,
			"publishedAt": t.PublishedAt,
		},
	})
}
