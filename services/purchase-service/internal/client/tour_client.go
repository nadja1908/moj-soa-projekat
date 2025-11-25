// purchase-service/internal/clients/tour_client.go
package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"purchase-service/internal/model"
)

type TourClient struct {
	client  *http.Client
	baseURL string
}

func NewTourClient() *TourClient {
	// Uzimanje URL-a iz Docker Compose env varijable
	baseURL := os.Getenv("TOUR_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://tour-service:8004" // Fallback za development
	}

	return &TourClient{
		client:  &http.Client{Timeout: 5 * time.Second},
		baseURL: baseURL,
	}
}

// GetPublishedTourDetails dohvaća ime, cenu i status objavljene ture
func (c *TourClient) GetPublishedTourDetails(tourID int64) (*model.TourPriceResponse, error) {
	url := fmt.Sprintf("%s/api/tours/public/%d", c.baseURL, tourID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Tour Service at %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("tour ID %d not found or not published", tourID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tour service error, status %d", resp.StatusCode)
	}

	var wrapper struct {
		Success bool                    `json:"success"`
		Tour    model.TourPriceResponse `json:"tour"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("failed to decode tour details: %w", err)
	}

	details := wrapper.Tour

	if details.Status != model.StatusPublished {
		return nil, fmt.Errorf("tour ID %d is not published, status: %s", tourID, details.Status)
	}

	return &details, nil
}
