package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"tour-service/internal/model"
	"tour-service/internal/store"

	"github.com/gin-gonic/gin"
	//"gorm.io/datatypes"
)

// ReviewHandler upravlja HTTP zahtevima vezanim za recenzije.
type ReviewHandler struct {
	Store store.ReviewStore
}

// NewReviewHandler kreira novu instancu ReviewHandler.
func NewReviewHandler(store store.ReviewStore) *ReviewHandler {
	return &ReviewHandler{Store: store}
}

// CreateReview omogućava turisti da ostavi recenziju za turu.
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	// 1. Dohvati userId iz konteksta (postavlja ga auth middleware)
	userIDInterface, exists := c.Get("userId")
	if !exists {
		log.Println("ERROR: userId not found in context for CreateReview")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated or ID missing"})
		return
	}

	touristID, ok := userIDInterface.(int64)
	if !ok {
		log.Println("ERROR: userId in context is not int64")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in context"})
		return
	}

	// 2. Parsiraj JSON body u DTO
	var req model.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 3. Validacija datuma posete ture (YYYY-MM-DD)
	dateVisited, err := time.Parse("2006-01-02", req.DateVisited)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dateVisited format. Expected YYYY-MM-DD"})
		return
	}

	// 4. Napravi Review model za bazu
	review := model.Review{
		TourID:      req.TourID,
		TouristID:   touristID,
		Rating:      req.Rating,
		Comment:     req.Comment,
		DateVisited: dateVisited,
	}

	// 5. Serijalizuj slike u JSON (čuvamo kao JSON u MySQL-u)
	// 5. Serijalizuj slike u JSON (čuvamo kao JSON u MySQL-u)
	if len(req.ImageURLs) > 0 {
		imageURLsBytes, err := json.Marshal(req.ImageURLs)
		if err != nil {
			log.Printf("WARN: Failed to serialize image URLs for review: %v", err)
			review.ImageURLs = []byte("[]")
		} else {
			review.ImageURLs = imageURLsBytes
		}
	} else {
		review.ImageURLs = []byte("[]")
	}

	// 6. Sačuvaj recenziju u bazi
	if err := h.Store.CreateReview(&review); err != nil {
		log.Printf("ERROR: Failed to create review: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save review"})
		return
	}

	// 7. Vrati uspešan odgovor
	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"message":  "Review created successfully",
		"reviewId": review.ID,
	})
}

// GetReviewsByTourID vraća sve recenzije za određenu turu (public).
func (h *ReviewHandler) GetReviewsByTourID(c *gin.Context) {
	// 1. Tour ID iz URL parametra
	tourIDStr := c.Param("tourId")
	tourID, err := strconv.ParseInt(tourIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tour ID format"})
		return
	}

	// 2. Učitaj recenzije iz baze
	reviews, err := h.Store.GetReviewsByTourID(tourID)
	if err != nil {
		log.Printf("ERROR: Failed to fetch reviews for tour %d: %v", tourID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	// 3. Mapiraj u DTO koji se vraća klijentu
	responseReviews := make([]model.ReviewResponse, 0, len(reviews))

	for _, review := range reviews {
		// --- imageUrls: datatypes.JSON ( []byte ) -> []string ---
		var images []string
		if review.ImageURLs != nil && len(review.ImageURLs) > 0 {
			if err := json.Unmarshal(review.ImageURLs, &images); err != nil {
				log.Printf("WARN: Failed to unmarshal image URLs for review %d: %v", review.ID, err)
				images = []string{}
			}
		} else {
			images = []string{}
		}

		// Za sada placeholder username/avatar
		touristUsername := "Tourist " + strconv.FormatInt(review.TouristID, 10)
		touristAvatar := "https://placehold.co/40x40/007bff/ffffff?text=U"

		responseReviews = append(responseReviews, model.ReviewResponse{
			ID:              review.ID,
			TourID:          review.TourID,
			TouristID:       review.TouristID,
			Rating:          review.Rating,
			Comment:         review.Comment,
			DateVisited:     review.DateVisited,
			DateCreated:     review.CreatedAt,
			TouristUsername: touristUsername,
			TouristAvatar:   touristAvatar,
			ImageURLs:       images,
		})
	}

	// 4. Pošalji odgovor
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"reviews": responseReviews,
	})
}
