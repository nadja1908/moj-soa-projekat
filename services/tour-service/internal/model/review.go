package model

import "time"

// Review predstavlja model recenzije u bazi podataka.
// Bez GORMA, čist struct za rad sa database/sql.
type Review struct {
	ID          uint      `json:"id"`
	TourID      int64     `json:"tourId"`
	TouristID   int64     `json:"touristId"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	DateVisited time.Time `json:"dateVisited"`
	CreatedAt   time.Time `json:"dateCreated"` // ručno dodato, umesto gorm.Model.CreatedAt
	ImageURLs   []byte    `json:"imageUrls"`   // čuvamo JSON kao []byte
}

// ReviewResponse je DTO za slanje recenzije klijentu.
type ReviewResponse struct {
	ID              uint      `json:"id"`
	TourID          int64     `json:"tourId"`
	TouristID       int64     `json:"touristId"`
	Rating          int       `json:"rating"`
	Comment         string    `json:"comment"`
	DateVisited     time.Time `json:"dateVisited"`
	DateCreated     time.Time `json:"dateCreated"`
	TouristUsername string    `json:"touristUsername"`
	TouristAvatar   string    `json:"touristAvatar"`
	ImageURLs       []string  `json:"imageUrls"`
}

// CreateReviewRequest – DTO koji primaš iz Postmana/frontenda.
type CreateReviewRequest struct {
	TourID      int64    `json:"tourId" binding:"required"`
	Rating      int      `json:"rating" binding:"required,min=1,max=5"`
	Comment     string   `json:"comment" binding:"required,max=500"`
	DateVisited string   `json:"dateVisited" binding:"required"` // format: YYYY-MM-DD
	ImageURLs   []string `json:"imageUrls"`
}
