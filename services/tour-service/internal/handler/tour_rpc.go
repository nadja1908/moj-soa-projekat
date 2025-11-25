package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"tour-service/internal/store"
	tourpb "tour-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TourRPCServer struct {
	tourpb.UnimplementedTourServiceServer
	store *store.Store
}

func NewTourRPCServer(store *store.Store) *TourRPCServer {
	return &TourRPCServer{
		store: store,
	}
}

// RPC: GetPublishedTours
func (s *TourRPCServer) GetPublishedTours(ctx context.Context, req *tourpb.GetPublishedToursRequest) (*tourpb.GetPublishedToursResponse, error) {
	log.Println("[RPC] GetPublishedTours called")

	tours, err := s.store.GetPublishedTours()
	if err != nil {
		log.Printf("[RPC] GetPublishedTours error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get published tours")
	}

	resp := &tourpb.GetPublishedToursResponse{}

	for _, t := range tours {
		// tags u bazi su JSON string -> []string
		var tags []string
		if t.Tags != "" {
			_ = json.Unmarshal([]byte(t.Tags), &tags)
		}

		publishedAt := ""
		if t.PublishedAt != nil {
			publishedAt = *t.PublishedAt
		}

		item := &tourpb.TourListItem{
			Id:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Difficulty:  string(t.Difficulty),
			Status:      string(t.Status),
			Price:       t.Price,
			DistanceKm:  t.DistanceKm,
			PublishedAt: publishedAt,
			Tags:        tags,
		}

		resp.Tours = append(resp.Tours, item)
	}

	return resp, nil
}

// RPC: GetPublicTour
func (s *TourRPCServer) GetPublicTour(ctx context.Context, req *tourpb.GetPublicTourRequest) (*tourpb.GetPublicTourResponse, error) {
	log.Printf("[RPC] GetPublicTour called for id=%d", req.Id)

	twd, err := s.store.GetTourForTourist(req.Id)
	if err != nil {
		log.Printf("[RPC] GetPublicTour error: %v", err)
		if err.Error() == "published tour not found" {
			return nil, status.Errorf(codes.NotFound, "tour not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get tour")
	}

	// tags iz JSON string -> []string
	var tags []string
	if twd.Tour.Tags != "" {
		_ = json.Unmarshal([]byte(twd.Tour.Tags), &tags)
	}

	// publishedAt u string
	var publishedAt string
	if twd.Tour.PublishedAt != nil && !twd.Tour.PublishedAt.IsZero() {
		// isti format kao u JSON: RFC3339
		publishedAt = twd.Tour.PublishedAt.Format(time.RFC3339)
	}

	// mapiraj durations
	var durations []*tourpb.TourDuration
	for _, d := range twd.Durations {
		durations = append(durations, &tourpb.TourDuration{
			TransportType:   string(d.TransportType),
			DurationMinutes: int32(d.DurationMinutes),
		})
	}

	// mapiraj keypoints (mi vraćamo samo prvu, ali kod podržava više)
	var keyPoints []*tourpb.KeyPoint
	for _, kp := range twd.KeyPoints {
		imageURL := ""
		if kp.ImageURL != nil {
			imageURL = *kp.ImageURL
		}
		keyPoints = append(keyPoints, &tourpb.KeyPoint{
			Name:        kp.Name,
			Description: kp.Description,
			Latitude:    kp.Latitude,
			Longitude:   kp.Longitude,
			ImageUrl:    imageURL,
		})
	}

	pbTour := &tourpb.PublicTour{
		Id:          twd.Tour.ID,
		Name:        twd.Tour.Name,
		Description: twd.Tour.Description,
		Difficulty:  string(twd.Tour.Difficulty),
		Status:      string(twd.Tour.Status),
		Price:       twd.Tour.Price,
		DistanceKm:  twd.Tour.DistanceKm,
		PublishedAt: publishedAt,
		Tags:        tags,
		Durations:   durations,
		KeyPoints:   keyPoints,
	}

	return &tourpb.GetPublicTourResponse{
		Tour: pbTour,
	}, nil
}
