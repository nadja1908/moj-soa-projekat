package rpc

import (
	"tour-service/internal/model"
)

// TourService RPC interface
type TourService struct {
	handler TourHandlerInterface
}

type TourHandlerInterface interface {
	ProcessCreateTour(userID int64, req *model.CreateTourRequest) (*model.Tour, error)
	ProcessGetTours(userID int64) ([]*model.Tour, error)
}

// NewTourService creates a new RPC tour service
func NewTourService(handler TourHandlerInterface) *TourService {
	return &TourService{handler: handler}
}

// CreateTour RPC method
func (s *TourService) CreateTour(args *CreateTourArgs, resp *model.Tour) error {
	result, err := s.handler.ProcessCreateTour(args.UserID, args.Request)
	if err != nil {
		return err
	}
	*resp = *result
	return nil
}

// GetTours RPC method
func (s *TourService) GetTours(userID *int64, resp *GetToursResponse) error {
	tours, err := s.handler.ProcessGetTours(*userID)
	if err != nil {
		return err
	}
	resp.Tours = tours
	return nil
}

// RPC argument types
type CreateTourArgs struct {
	UserID  int64                    `json:"userId"`
	Request *model.CreateTourRequest `json:"request"`
}

type GetToursResponse struct {
	Tours []*model.Tour `json:"tours"`
}
