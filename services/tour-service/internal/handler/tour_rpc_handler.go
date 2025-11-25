package handler

import (
	"context"
	"tour-service/internal/store"
	tourpb "tour-service/proto"
)

type TourRPCHandler struct {
	tourpb.UnimplementedTourServiceServer
	store *store.Store
}

func NewTourRPCHandler(store *store.Store) *TourRPCHandler {
	return &TourRPCHandler{store: store}
}

func (h *TourRPCHandler) GetPublishedTours(ctx context.Context, req *tourpb.GetPublishedToursRequest) (*tourpb.GetPublishedToursResponse, error) {
	tours, err := h.store.GetPublishedTours()
	if err != nil {
		return nil, err
	}

	resp := &tourpb.GetPublishedToursResponse{}
	for _, t := range tours {
		resp.Tours = append(resp.Tours, &tourpb.TourListItem{
			Id:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Difficulty:  string(t.Difficulty),
			Status:      string(t.Status),
			Price:       t.Price,
			DistanceKm:  t.DistanceKm,
			PublishedAt: safeStrPtr(t.PublishedAt),
			Tags:        []string{}, // možeš da dodaš kasnije
		})
	}

	return resp, nil
}

func safeStrPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
