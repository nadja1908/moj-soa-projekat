package model

import (
	"database/sql"
	"time"
)

// TourStatus predstavlja status ture
type TourStatus string

const (
	StatusDraft     TourStatus = "DRAFT"
	StatusPublished TourStatus = "PUBLISHED"
	StatusArchived  TourStatus = "ARCHIVED"
)

// TourDifficulty predstavlja težinu ture
type TourDifficulty string

const (
	DifficultyEasy     TourDifficulty = "EASY"
	DifficultyModerate TourDifficulty = "MODERATE"
	DifficultyHard     TourDifficulty = "HARD"
)

// TransportType predstavlja tip prevoza
type TransportType string

const (
	TransportWalk    TransportType = "walk"
	TransportBicycle TransportType = "bicycle"
	TransportCar     TransportType = "car"
)

// Tour predstavlja turu
type Tour struct {
	ID          int64          `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	Difficulty  TourDifficulty `json:"difficulty" db:"difficulty_level"`
	Tags        string         `json:"tags" db:"tags"` // JSON array kao string
	Status      TourStatus     `json:"status" db:"status"`
	Price       float64        `json:"price" db:"price"`
	AuthorID    int64          `json:"authorId" db:"author_id"`
	DistanceKm  float64        `json:"distanceKm" db:"distance_km"`
	PublishedAt *time.Time     `json:"publishedAt" db:"published_at"`
	ArchivedAt  *time.Time     `json:"archivedAt" db:"archived_at"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time      `json:"updatedAt" db:"updated_at"`
}

// KeyPoint predstavlja ključnu tačku ture
type KeyPoint struct {
	ID          int64     `json:"id" db:"id"`
	TourID      int64     `json:"tourId" db:"tour_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Latitude    float64   `json:"latitude" db:"latitude"`
	Longitude   float64   `json:"longitude" db:"longitude"`
	OrderIndex  int       `json:"orderIndex" db:"order_index"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// TourDuration predstavlja vreme potrebno za obiđanje ture po različitim vrstama prevoza
type TourDuration struct {
	ID              int64         `json:"id" db:"id"`
	TourID          int64         `json:"tourId" db:"tour_id"`
	TransportType   TransportType `json:"transportType" db:"transport_type"`
	DurationMinutes int           `json:"durationMinutes" db:"duration_minutes"`
	CreatedAt       time.Time     `json:"createdAt" db:"created_at"`
}

// TourExecution predstavlja izvršavanje ture od strane turiste
type TourExecution struct {
	ID               int64      `json:"id" db:"id"`
	TourID           int64      `json:"tourId" db:"tour_id"`
	TouristID        int64      `json:"touristId" db:"tourist_id"`
	CurrentLatitude  *float64   `json:"currentLatitude" db:"current_latitude"`
	CurrentLongitude *float64   `json:"currentLongitude" db:"current_longitude"`
	StartedAt        time.Time  `json:"startedAt" db:"started_at"`
	CompletedAt      *time.Time `json:"completedAt" db:"completed_at"`
}

// TourWithDetails predstavlja turu sa svim povezanim podacima
type TourWithDetails struct {
	Tour      Tour           `json:"tour"`
	KeyPoints []KeyPoint     `json:"keyPoints"`
	Durations []TourDuration `json:"durations"`
}

// TourListItem predstavlja turu u listi (bez svih detalja)
type TourListItem struct {
	ID            int64          `json:"id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Difficulty    TourDifficulty `json:"difficulty"`
	Status        TourStatus     `json:"status"`
	Price         float64        `json:"price"`
	DistanceKm    float64        `json:"distanceKm"`
	PublishedAt   sql.NullTime   `json:"publishedAt"`
	FirstKeyPoint *KeyPoint      `json:"firstKeyPoint"`
}
