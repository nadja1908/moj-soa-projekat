package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourStatus string

const (
	StatusPublished TourStatus = "PUBLISHED"
)

// --- Entiteti za MongoDB ---

// OrderItem - Stavka u korpi
type OrderItem struct {
	TourID   int64   `json:"tourId" bson:"tourId"`
	TourName string  `json:"tourName" bson:"tourName"`
	Price    float64 `json:"price" bson:"price"`
	Quantity int     `json:"quantity" bson:"quantity"`
}

// ShoppingCart - Korpa (vezana za jednog Turistu)
type ShoppingCart struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TouristID   int64              `json:"touristId" bson:"touristId"`
	Items       []OrderItem        `json:"items" bson:"items"`
	TotalPrice  float64            `json:"totalPrice" bson:"totalPrice"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	LastUpdated time.Time          `json:"lastUpdated" bson:"lastUpdated"`
}

// TourPurchaseToken - Token koji se dobija nakon kupovine
type TourPurchaseToken struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token        string             `json:"token" bson:"token"` // UUID
	TouristID    int64              `json:"touristId" bson:"touristId"`
	TourID       int64              `json:"tourId" bson:"tourId"`
	PricePaid    float64            `json:"pricePaid" bson:"pricePaid"`
	PurchaseDate time.Time          `json:"purchaseDate" bson:"purchaseDate"`
}

// --- DTOs (Data Transfer Objects) ---

// AddToCartRequest - Zahtev od klijenta
type AddToCartRequest struct {
	TourID   int64 `json:"tourId" binding:"required"`
	Quantity int   `json:"quantity"` // Podrazumevana vrednost 1
}

// TourPriceResponse - Odgovor od Tour Service-a
type TourPriceResponse struct {
	ID     int64      `json:"id"`
	Name   string     `json:"name"`
	Price  float64    `json:"price"`
	Status TourStatus `json:"status"`
}
