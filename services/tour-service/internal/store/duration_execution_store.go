package store

import (
	"fmt"

	"tour-service/internal/model"
)

// CreateTourDuration kreira novo vreme trajanja ture
func (s *Store) CreateTourDuration(req *model.CreateTourDurationRequest) (*model.TourDuration, error) {
	query := `
		INSERT INTO tour_durations (tour_id, transport_type, duration_minutes)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		duration_minutes = VALUES(duration_minutes)
	`

	result, err := s.db.Exec(query, req.TourID, req.TransportType, req.DurationMinutes)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		// Ako je update, dobij postojeći ID
		return s.GetTourDurationByTourAndTransport(req.TourID, req.TransportType)
	}

	return s.GetTourDurationByID(id)
}

// GetTourDurationByID vraća vreme trajanja po ID-u
func (s *Store) GetTourDurationByID(id int64) (*model.TourDuration, error) {
	query := `
		SELECT id, tour_id, transport_type, duration_minutes, created_at
		FROM tour_durations WHERE id = ?
	`

	row := s.db.QueryRow(query, id)
	duration := &model.TourDuration{}

	err := row.Scan(
		&duration.ID, &duration.TourID, &duration.TransportType,
		&duration.DurationMinutes, &duration.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return duration, nil
}

// GetTourDurationByTourAndTransport vraća vreme trajanja po tour_id i transport_type
func (s *Store) GetTourDurationByTourAndTransport(tourID int64, transportType model.TransportType) (*model.TourDuration, error) {
	query := `
		SELECT id, tour_id, transport_type, duration_minutes, created_at
		FROM tour_durations WHERE tour_id = ? AND transport_type = ?
	`

	row := s.db.QueryRow(query, tourID, transportType)
	duration := &model.TourDuration{}

	err := row.Scan(
		&duration.ID, &duration.TourID, &duration.TransportType,
		&duration.DurationMinutes, &duration.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return duration, nil
}

// GetTourDurations vraća sva vremena trajanja za turu
func (s *Store) GetTourDurations(tourID int64) ([]model.TourDuration, error) {
	query := `
		SELECT id, tour_id, transport_type, duration_minutes, created_at
		FROM tour_durations 
		WHERE tour_id = ?
		ORDER BY transport_type
	`

	rows, err := s.db.Query(query, tourID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var durations []model.TourDuration
	for rows.Next() {
		duration := model.TourDuration{}
		err := rows.Scan(
			&duration.ID, &duration.TourID, &duration.TransportType,
			&duration.DurationMinutes, &duration.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		durations = append(durations, duration)
	}

	return durations, nil
}

// DeleteTourDuration briše vreme trajanja
func (s *Store) DeleteTourDuration(id int64) error {
	query := "DELETE FROM tour_durations WHERE id = ?"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tour duration not found")
	}

	return nil
}

// StartTourExecution počinje izvršavanje ture
func (s *Store) StartTourExecution(touristID, tourID int64) (*model.TourExecution, error) {
	query := `
		INSERT INTO tour_executions (tour_id, tourist_id)
		VALUES (?, ?)
	`

	result, err := s.db.Exec(query, tourID, touristID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetTourExecutionByID(id)
}

// GetTourExecutionByID vraća izvršavanje ture po ID-u
func (s *Store) GetTourExecutionByID(id int64) (*model.TourExecution, error) {
	query := `
		SELECT id, tour_id, tourist_id, current_latitude, current_longitude,
		       started_at, completed_at
		FROM tour_executions WHERE id = ?
	`

	row := s.db.QueryRow(query, id)
	execution := &model.TourExecution{}

	err := row.Scan(
		&execution.ID, &execution.TourID, &execution.TouristID,
		&execution.CurrentLatitude, &execution.CurrentLongitude,
		&execution.StartedAt, &execution.CompletedAt,
	)

	if err != nil {
		return nil, err
	}

	return execution, nil
}

// GetActiveTourExecution vraća aktivno izvršavanje ture za turista
func (s *Store) GetActiveTourExecution(touristID int64) (*model.TourExecution, error) {
	query := `
		SELECT id, tour_id, tourist_id, current_latitude, current_longitude,
		       started_at, completed_at
		FROM tour_executions 
		WHERE tourist_id = ? AND completed_at IS NULL
		ORDER BY started_at DESC
		LIMIT 1
	`

	row := s.db.QueryRow(query, touristID)
	execution := &model.TourExecution{}

	err := row.Scan(
		&execution.ID, &execution.TourID, &execution.TouristID,
		&execution.CurrentLatitude, &execution.CurrentLongitude,
		&execution.StartedAt, &execution.CompletedAt,
	)

	if err != nil {
		return nil, err
	}

	return execution, nil
}

// UpdateTouristLocation ažurira trenutnu lokaciju turiste
func (s *Store) UpdateTouristLocation(touristID int64, latitude, longitude float64) error {
	// Ažuriraj aktivno izvršavanje ture
	query := `
		UPDATE tour_executions 
		SET current_latitude = ?, current_longitude = ?
		WHERE tourist_id = ? AND completed_at IS NULL
	`

	result, err := s.db.Exec(query, latitude, longitude, touristID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no active tour execution found")
	}

	return nil
}

// CompleteTourExecution označava izvršavanje ture kao završeno
func (s *Store) CompleteTourExecution(executionID int64) error {
	query := `
		UPDATE tour_executions 
		SET completed_at = NOW()
		WHERE id = ? AND completed_at IS NULL
	`

	result, err := s.db.Exec(query, executionID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tour execution not found or already completed")
	}

	return nil
}

// GetTourWithDetails vraća kompletan tour sa svim povezanim podacima
func (s *Store) GetTourWithDetails(id int64) (*model.TourWithDetails, error) {
	tour, err := s.GetTourByID(id)
	if err != nil {
		return nil, err
	}

	keyPoints, err := s.GetKeyPoints(id)
	if err != nil {
		return nil, err
	}

	durations, err := s.GetTourDurations(id)
	if err != nil {
		return nil, err
	}

	return &model.TourWithDetails{
		Tour:      *tour,
		KeyPoints: keyPoints,
		Durations: durations,
	}, nil
}
