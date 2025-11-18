package model

// CreateTourRequest je model za kreiranje ture
type CreateTourRequest struct {
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Difficulty  TourDifficulty `json:"difficulty" binding:"required"`
	Tags        []string       `json:"tags"`
}

// UpdateTourRequest je model za ažuriranje ture
type UpdateTourRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Difficulty  TourDifficulty `json:"difficulty"`
	Tags        []string       `json:"tags"`
	Price       *float64       `json:"price"`
}

// CreateKeyPointRequest je model za kreiranje ključne tačke
type CreateKeyPointRequest struct {
	TourID      int64   `json:"tourId" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
}

// UpdateKeyPointRequest je model za ažuriranje ključne tačke
type UpdateKeyPointRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
}

// CreateTourDurationRequest je model za kreiranje vremena trajanja ture
type CreateTourDurationRequest struct {
	TourID          int64         `json:"tourId" binding:"required"`
	TransportType   TransportType `json:"transportType" binding:"required"`
	DurationMinutes int           `json:"durationMinutes" binding:"required"`
}

// PublishTourRequest je model za objavljivanje ture
type PublishTourRequest struct {
	Price float64 `json:"price" binding:"required,min=0"`
}

// UpdateLocationRequest je model za ažuriranje lokacije turiste
type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// StartTourExecutionRequest je model za početak izvršavanja ture
type StartTourExecutionRequest struct {
	TourID int64 `json:"tourId" binding:"required"`
}

// TourResponse je odgovor sa osnovnim podacima o turi
type TourResponse struct {
	Success bool             `json:"success"`
	Tour    *TourWithDetails `json:"tour,omitempty"`
	Error   string           `json:"error,omitempty"`
}

// ToursListResponse je odgovor sa listom tura
type ToursListResponse struct {
	Success bool           `json:"success"`
	Tours   []TourListItem `json:"tours,omitempty"`
	Error   string         `json:"error,omitempty"`
}

// KeyPointResponse je odgovor sa podacima o ključnoj tački
type KeyPointResponse struct {
	Success  bool      `json:"success"`
	KeyPoint *KeyPoint `json:"keyPoint,omitempty"`
	Error    string    `json:"error,omitempty"`
}
