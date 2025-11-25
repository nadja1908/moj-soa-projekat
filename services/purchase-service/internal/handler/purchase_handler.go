// purchase-service/internal/handler/purchase_handler.go
package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"purchase-service/internal/client"
	"strconv"
	"time"

	"purchase-service/internal/model"
	"purchase-service/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PurchaseHandler struct {
	store      *store.Store
	tourClient *client.TourClient
}

func NewPurchaseHandler(store *store.Store, tourClient *client.TourClient) *PurchaseHandler {
	return &PurchaseHandler{store: store, tourClient: tourClient}
}

func getTouristIDFromContext(c *gin.Context) int64 {
	userIDValue, exists := c.Get("userID")
	if !exists {
		return 0
	}
	userID, ok := userIDValue.(int64)
	if !ok {
		return 0
	}
	return userID
}

func (h *PurchaseHandler) AddToCart(c *gin.Context) {
	touristID := getTouristIDFromContext(c)
	if touristID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid Tourist ID"})
		return
	}

	var req model.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	if req.Quantity == 0 {
		req.Quantity = 1
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	tourDetails, err := h.tourClient.GetPublishedTourDetails(req.TourID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add tour: " + err.Error()})
		return
	}

	newItem := model.OrderItem{
		TourID:   tourDetails.ID,
		TourName: tourDetails.Name,
		Price:    tourDetails.Price,
		Quantity: req.Quantity,
	}

	updatedCart, err := h.store.SaveCartWithItem(ctx, touristID, newItem)
	if err != nil {
		log.Printf("Error saving cart: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shopping cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tour added to cart successfully",
		"cart":    updatedCart,
	})
}

func (h *PurchaseHandler) GetShoppingCart(c *gin.Context) {
	touristID := getTouristIDFromContext(c)
	if touristID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid Tourist ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cart, err := h.store.GetCartByTouristID(ctx, touristID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shopping cart"})
		return
	}

	if cart == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Shopping cart is empty", "cart": model.ShoppingCart{TouristID: touristID}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func (h *PurchaseHandler) RemoveItem(c *gin.Context) {
	touristID := getTouristIDFromContext(c)
	if touristID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid Tourist ID"})
		return
	}

	tourIDStr := c.Param("tourId")
	tourID, err := strconv.ParseInt(tourIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tour ID in path"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedCart, err := h.store.RemoveItemFromCart(ctx, touristID, tourID)
	if err != nil {
		if err.Error() == "shopping cart not found" || err.Error() == fmt.Sprintf("item with tour ID %d not found in cart", tourID) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed successfully",
		"cart":    updatedCart,
	})
}

func (h *PurchaseHandler) Checkout(c *gin.Context) {
	touristID := getTouristIDFromContext(c)
	if touristID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid Tourist ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	cart, err := h.store.GetCartByTouristID(ctx, touristID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shopping cart"})
		return
	}
	if cart == nil || len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Shopping cart is empty"})
		return
	}

	var generatedTokens []model.TourPurchaseToken
	for _, item := range cart.Items {
		for i := 0; i < item.Quantity; i++ {
			token := model.TourPurchaseToken{
				Token:        uuid.New().String(),
				TouristID:    touristID,
				TourID:       item.TourID,
				PricePaid:    item.Price,
				PurchaseDate: time.Now(),
			}
			generatedTokens = append(generatedTokens, token)
		}
	}

	if err = h.store.SaveTokens(ctx, generatedTokens); err != nil {
		log.Printf("CRITICAL ERROR: Failed to save purchase tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment failed (Internal server error during token generation)"})
		return
	}

	if err = h.store.ClearCart(ctx, touristID); err != nil {
		log.Printf("Warning: Failed to clear cart after successful checkout for tourist %d: %v", touristID, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Successfully purchased %d tours.", len(generatedTokens)),
		"tokens":  generatedTokens,
	})
}

func (h *PurchaseHandler) InternalGetCart(c *gin.Context) {
	touristIDStr := c.Param("touristID")
	touristID, err := strconv.ParseInt(touristIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tourist ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cart, err := h.store.GetCartByTouristID(ctx, touristID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shopping cart"})
		return
	}

	if cart == nil {
		c.JSON(http.StatusOK, gin.H{"cart": model.ShoppingCart{TouristID: touristID}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func (h *PurchaseHandler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	err := h.store.Client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Health check failed: MongoDB connection error: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "service": "purchase-service", "db_status": "unavailable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "purchase-service", "db_status": "ok"})
}
