package handler

import (
	"net/http"
	"strconv"

	"tour-service/internal/model"
	"tour-service/internal/store"

	"github.com/gin-gonic/gin"
)

type TourDurationHandler struct {
	store *store.Store
}

func NewTourDurationHandler(store *store.Store) *TourDurationHandler {
	return &TourDurationHandler{store: store}
}

// CreateTourDuration kreira novo vreme trajanja ture
func (h *TourDurationHandler) CreateTourDuration(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req model.CreateTourDurationRequest
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
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only add durations to your own tours"})
		return
	}

	// Proverava da li je tura u draft stanju
	if tour.Status != model.StatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only add durations to draft tours"})
		return
	}

	duration, err := h.store.CreateTourDuration(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tour duration"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"duration": duration,
	})
}

// GetTourDurations vraća sva vremena trajanja za turu
func (h *TourDurationHandler) GetTourDurations(c *gin.Context) {
	tourID, err := strconv.ParseInt(c.Param("tourId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID"})
		return
	}

	durations, err := h.store.GetTourDurations(tourID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tour durations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"durations": durations,
	})
}

// DeleteTourDuration briše vreme trajanja
func (h *TourDurationHandler) DeleteTourDuration(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration ID"})
		return
	}

	duration, err := h.store.GetTourDurationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour duration not found"})
		return
	}

	// Proverava da li je korisnik vlasnik ture
	tour, err := h.store.GetTourByID(duration.TourID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	if tour.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete durations from your own tours"})
		return
	}

	// Proverava da li je tura u draft stanju
	if tour.Status != model.StatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only delete durations from draft tours"})
		return
	}

	err = h.store.DeleteTourDuration(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tour duration deleted successfully",
	})
}

type TourSimulatorHandler struct {
	store *store.Store
}

func NewTourSimulatorHandler(store *store.Store) *TourSimulatorHandler {
	return &TourSimulatorHandler{store: store}
}

// UpdateLocation ažurira trenutnu lokaciju turiste (simulator pozicije)
func (h *TourSimulatorHandler) UpdateLocation(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req model.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.store.UpdateTouristLocation(userID, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Location updated successfully",
		"latitude":  req.Latitude,
		"longitude": req.Longitude,
	})
}

// GetCurrentLocation vraća trenutnu lokaciju turiste
func (h *TourSimulatorHandler) GetCurrentLocation(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	execution, err := h.store.GetActiveTourExecution(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active tour execution found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"execution": execution,
	})
}

// StartTourExecution počinje izvršavanje ture
func (h *TourSimulatorHandler) StartTourExecution(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req model.StartTourExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Proverava da li je tura objavljena
	tour, err := h.store.GetTourByID(req.TourID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	if tour.Status != model.StatusPublished {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only start execution of published tours"})
		return
	}

	// Proverava da li već postoji aktivno izvršavanje
	_, err = h.store.GetActiveTourExecution(userID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already have an active tour execution"})
		return
	}

	execution, err := h.store.StartTourExecution(userID, req.TourID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start tour execution"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":   true,
		"execution": execution,
		"message":   "Tour execution started successfully",
	})
}

// CompleteTourExecution završava izvršavanje ture
func (h *TourSimulatorHandler) CompleteTourExecution(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid execution ID"})
		return
	}

	// Proverava da li je korisnik vlasnik izvršavanja
	execution, err := h.store.GetTourExecutionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour execution not found"})
		return
	}

	if execution.TouristID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only complete your own tour executions"})
		return
	}

	err = h.store.CompleteTourExecution(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tour execution completed successfully",
	})
}
