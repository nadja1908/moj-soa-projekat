package handler

import (
	"tour-service/internal/model"
	"tour-service/internal/store"
)

// TourRPCHandler implements the RPC interface for TourHandler
type TourRPCHandler struct {
	store *store.Store
}

func NewTourRPCHandler(store *store.Store) *TourRPCHandler {
	return &TourRPCHandler{store: store}
}

// ProcessCreateTour handles tour creation logic for RPC
func (h *TourRPCHandler) ProcessCreateTour(userID int64, req *model.CreateTourRequest) (*model.Tour, error) {
	if userID == 0 {
		return nil, &RPCError{Code: 401, Message: "User not authenticated"}
	}

	tour, err := h.store.CreateTour(userID, req)
	if err != nil {
		return nil, &RPCError{Code: 500, Message: "Failed to create tour"}
	}

	return tour, nil
}

// ProcessGetTours handles getting tours logic for RPC
func (h *TourRPCHandler) ProcessGetTours(userID int64) ([]*model.Tour, error) {
	if userID == 0 {
		return nil, &RPCError{Code: 401, Message: "User not authenticated"}
	}

	tourItems, err := h.store.GetToursByAuthor(userID)
	if err != nil {
		return nil, &RPCError{Code: 500, Message: "Failed to get tours"}
	}

	// Convert TourListItem to Tour pointers
	tours := make([]*model.Tour, len(tourItems))
	for i, item := range tourItems {
		tours[i] = &model.Tour{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			AuthorID:    userID, // Use the userID since TourListItem doesn't have AuthorID
			Difficulty:  item.Difficulty,
			Tags:        item.Tags,
			Status:      item.Status,
			DistanceKm:  item.DistanceKm,
			Price:       item.Price,
		}
	}

	return tours, nil
}

// RPCError represents an RPC error with status code
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RPCError) Error() string {
	return e.Message
}
