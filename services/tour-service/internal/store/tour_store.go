package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"tour-service/internal/model"
)

// CreateTour kreira novu turu
func (s *Store) CreateTour(authorID int64, req *model.CreateTourRequest) (*model.Tour, error) {
	tagsJSON, _ := json.Marshal(req.Tags)

	query := `
		INSERT INTO tours (name, description, difficulty_level, tags, author_id, status, price)
		VALUES (?, ?, ?, ?, ?, 'DRAFT', 0.00)
	`

	result, err := s.db.Exec(query, req.Name, req.Description, req.Difficulty, tagsJSON, authorID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetTourByID(id)
}

// GetTourByID vraća turu po ID-u
func (s *Store) GetTourByID(id int64) (*model.Tour, error) {
	query := `
		SELECT id, name, description, difficulty_level, tags, status, price, author_id,
		       distance_km, published_at, archived_at, created_at, updated_at
		FROM tours WHERE id = ?
	`

	row := s.db.QueryRow(query, id)
	tour := &model.Tour{}
	var tagsJSON sql.NullString

	err := row.Scan(
		&tour.ID, &tour.Name, &tour.Description, &tour.Difficulty, &tagsJSON,
		&tour.Status, &tour.Price, &tour.AuthorID, &tour.DistanceKm,
		&tour.PublishedAt, &tour.ArchivedAt, &tour.CreatedAt, &tour.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Parse tags JSON
	if tagsJSON.Valid {
		tour.Tags = tagsJSON.String
	} else {
		tour.Tags = "[]"
	}

	return tour, nil
}

// GetToursByAuthor vraća sve ture određenog autora
func (s *Store) GetToursByAuthor(authorID int64) ([]model.TourListItem, error) {
	query := `
		SELECT t.id, t.name, t.description, t.difficulty_level, t.status, t.price,
		       t.distance_km, DATE_FORMAT(t.published_at, '%Y-%m-%d %H:%i:%s') as published_at, t.tags, t.created_at
		FROM tours t
		WHERE t.author_id = ?
		ORDER BY t.created_at DESC
	`

	rows, err := s.db.Query(query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tours []model.TourListItem
	for rows.Next() {
		tour := model.TourListItem{}
		var tagsJSON sql.NullString
		var publishedAtStr sql.NullString

		err := rows.Scan(
			&tour.ID, &tour.Name, &tour.Description, &tour.Difficulty,
			&tour.Status, &tour.Price, &tour.DistanceKm, &publishedAtStr,
			&tagsJSON, &tour.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert sql.NullString datetime to *string (samo datum)
		if publishedAtStr.Valid && publishedAtStr.String != "" {
			if parsedTime, err := time.Parse("2006-01-02 15:04:05", publishedAtStr.String); err == nil {
				dateOnly := parsedTime.Format("2006-01-02") // Samo datum do 'T'
				tour.PublishedAt = &dateOnly
			} else {
				tour.PublishedAt = nil
			}
		} else {
			tour.PublishedAt = nil
		}

		// Parse tags JSON
		if tagsJSON.Valid {
			tour.Tags = tagsJSON.String
		} else {
			tour.Tags = "[]"
		}

		// Dodaj samo prvu ključnu tačku (ostale se ne pokazuju turistima)
		// firstKeyPoint, _ := s.GetFirstKeyPoint(tour.ID)
		// tour.FirstKeyPoint = firstKeyPoint
		tour.FirstKeyPoint = nil // Temporarily disabled

		tours = append(tours, tour)
	}

	return tours, nil
}

// GetPublishedTours vraća sve objavljene ture (za turiste)
func (s *Store) GetPublishedTours() ([]model.TourListItem, error) {
	query := `
		SELECT t.id, t.name, t.description, t.difficulty_level, t.status, t.price,
		       t.distance_km, DATE_FORMAT(t.published_at, '%Y-%m-%d %H:%i:%s') as published_at, t.tags, t.created_at
		FROM tours t
		WHERE t.status = 'PUBLISHED'
		ORDER BY t.published_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tours []model.TourListItem
	for rows.Next() {
		tour := model.TourListItem{}
		var tagsJSON sql.NullString
		var publishedAtStr sql.NullString

		err := rows.Scan(
			&tour.ID, &tour.Name, &tour.Description, &tour.Difficulty,
			&tour.Status, &tour.Price, &tour.DistanceKm, &publishedAtStr,
			&tagsJSON, &tour.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert sql.NullString datetime to *string (samo datum)
		if publishedAtStr.Valid && publishedAtStr.String != "" {
			if parsedTime, err := time.Parse("2006-01-02 15:04:05", publishedAtStr.String); err == nil {
				dateOnly := parsedTime.Format("2006-01-02") // Samo datum do 'T'
				tour.PublishedAt = &dateOnly
			} else {
				tour.PublishedAt = nil
			}
		} else {
			tour.PublishedAt = nil
		}

		// Parse tags JSON
		if tagsJSON.Valid {
			tour.Tags = tagsJSON.String
		} else {
			tour.Tags = "[]"
		}

		// Dodaj samo prvu ključnu tačku (ostale se ne pokazuju turistima)
		firstKeyPoint, _ := s.GetFirstKeyPoint(tour.ID)
		tour.FirstKeyPoint = firstKeyPoint

		tours = append(tours, tour)
	}

	return tours, nil
}

// UpdateTour ažurira turu
func (s *Store) UpdateTour(id int64, req *model.UpdateTourRequest) (*model.Tour, error) {
	setParts := []string{}
	args := []interface{}{}

	if req.Name != "" {
		setParts = append(setParts, "name = ?")
		args = append(args, req.Name)
	}
	if req.Description != "" {
		setParts = append(setParts, "description = ?")
		args = append(args, req.Description)
	}
	if req.Difficulty != "" {
		setParts = append(setParts, "difficulty = ?")
		args = append(args, req.Difficulty)
	}
	if req.Tags != nil {
		tagsJSON, _ := json.Marshal(req.Tags)
		setParts = append(setParts, "tags = ?")
		args = append(args, string(tagsJSON))
	}
	if req.Price != nil {
		setParts = append(setParts, "price = ?")
		args = append(args, *req.Price)
	}

	if len(setParts) == 0 {
		return s.GetTourByID(id)
	}

	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE tours SET %s WHERE id = ?", fmt.Sprintf("%s", setParts[0]))
	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf("%s, %s", query, setParts[i])
	}

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetTourByID(id)
}

// PublishTour objavljuje turu
func (s *Store) PublishTour(id int64, price float64) error {
	// Proveri uslove za objavljivanje
	tour, err := s.GetTourByID(id)
	if err != nil {
		return err
	}

	if tour.Status != model.StatusDraft {
		return fmt.Errorf("tour must be in draft status to publish")
	}

	// Proveri da li ima bar 2 ključne tačke
	keyPointsCount, err := s.GetKeyPointsCount(id)
	if err != nil {
		return err
	}
	if keyPointsCount < 2 {
		return fmt.Errorf("tour must have at least 2 key points")
	}

	// Proveri da li ima bar jedno vreme trajanja
	durationsCount, err := s.GetDurationsCount(id)
	if err != nil {
		return err
	}
	if durationsCount < 1 {
		return fmt.Errorf("tour must have at least one duration defined")
	}

	query := `
		UPDATE tours 
		SET status = 'PUBLISHED', price = ?, published_at = NOW(), updated_at = NOW()
		WHERE id = ?
	`

	_, err = s.db.Exec(query, price, id)
	return err
}

// ArchiveTour arhivira turu
func (s *Store) ArchiveTour(id int64) error {
	query := `
		UPDATE tours 
		SET status = 'ARCHIVED', archived_at = NOW(), updated_at = NOW()
		WHERE id = ? AND status = 'PUBLISHED'
	`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tour not found or not published")
	}

	return nil
}

// ReactivateTour ponovo aktivira arhivirtanu turu
func (s *Store) ReactivateTour(id int64) error {
	query := `
		UPDATE tours 
		SET status = 'PUBLISHED', archived_at = NULL, updated_at = NOW()
		WHERE id = ? AND status = 'ARCHIVED'
	`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tour not found or not archived")
	}

	return nil
}

// DeleteTour briše turu (samo draft ture)
func (s *Store) DeleteTour(id int64, authorID int64) error {
	query := `DELETE FROM tours WHERE id = ? AND author_id = ? AND status = 'draft'`

	result, err := s.db.Exec(query, id, authorID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tour not found or cannot be deleted")
	}

	return nil
}

// GetKeyPointsCount vraća broj ključnih tačaka za turu
func (s *Store) GetKeyPointsCount(tourID int64) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM key_points WHERE tour_id = ?"
	err := s.db.QueryRow(query, tourID).Scan(&count)
	return count, err
}

// GetDurationsCount vraća broj definisanih trajanja za turu
func (s *Store) GetDurationsCount(tourID int64) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM tour_durations WHERE tour_id = ?"
	err := s.db.QueryRow(query, tourID).Scan(&count)
	return count, err
}

// GetFirstKeyPoint vraća prvu ključnu tačku ture
func (s *Store) GetFirstKeyPoint(tourID int64) (*model.KeyPoint, error) {
	query := `
		SELECT id, tour_id, name, description, latitude, longitude,
		       order_index, created_at, updated_at
		FROM key_points 
		WHERE tour_id = ? 
		ORDER BY order_index ASC 
		LIMIT 1
	`

	row := s.db.QueryRow(query, tourID)
	keyPoint := &model.KeyPoint{}

	err := row.Scan(
		&keyPoint.ID, &keyPoint.TourID, &keyPoint.Name, &keyPoint.Description,
		&keyPoint.Latitude, &keyPoint.Longitude,
		&keyPoint.OrderIndex, &keyPoint.CreatedAt, &keyPoint.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return keyPoint, nil
}

// GetTourForTourist vraća objavljenu turu sa osnovnim podacima i samo prvom ključnom tačkom
func (s *Store) GetTourForTourist(id int64) (*model.TourWithDetails, error) {
	// Prvo proveri da li je tura objavljena
	query := `
		SELECT id, name, description, difficulty_level, tags, status, price, author_id,
		       distance_km, published_at, archived_at, created_at, updated_at
		FROM tours WHERE id = ? AND status = 'PUBLISHED'
	`

	row := s.db.QueryRow(query, id)
	tour := &model.Tour{}
	var tagsJSON sql.NullString

	err := row.Scan(
		&tour.ID, &tour.Name, &tour.Description, &tour.Difficulty, &tagsJSON,
		&tour.Status, &tour.Price, &tour.AuthorID, &tour.DistanceKm,
		&tour.PublishedAt, &tour.ArchivedAt, &tour.CreatedAt, &tour.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("published tour not found")
	}
	if err != nil {
		return nil, err
	}

	// Parse tags
	if tagsJSON.Valid {
		tour.Tags = tagsJSON.String
	}

	// Dobij samo prvu ključnu tačku
	firstKeyPoint, err := s.GetFirstKeyPoint(id)
	var keyPoints []model.KeyPoint
	if err == nil && firstKeyPoint != nil {
		keyPoints = []model.KeyPoint{*firstKeyPoint}
		fmt.Printf("DEBUG: GetTourForTourist - Found first key point: %+v\n", *firstKeyPoint)
	} else {
		fmt.Printf("DEBUG: GetTourForTourist - No first key point found, err: %v\n", err)
	}

	// Dobij durations
	durations, err := s.GetTourDurations(id)
	if err != nil {
		fmt.Printf("DEBUG: GetTourForTourist - Error getting durations: %v\n", err)
		durations = []model.TourDuration{}
	} else {
		fmt.Printf("DEBUG: GetTourForTourist - Found %d durations\n", len(durations))
	}

	return &model.TourWithDetails{
		Tour:      *tour,
		KeyPoints: keyPoints,
		Durations: durations,
	}, nil
}

// CalculateDistance računa udaljenost između dve geografske tačke u kilometrima
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // km

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// UpdateTourDistance ažurira ukupnu dužinu ture na osnovu ključnih tačaka
func (s *Store) UpdateTourDistance(tourID int64) error {
	keyPoints, err := s.GetKeyPoints(tourID)
	if err != nil {
		return err
	}

	if len(keyPoints) < 2 {
		return nil // Nema dovoljno tačaka za računanje distance
	}

	var totalDistance float64
	for i := 1; i < len(keyPoints); i++ {
		prev := keyPoints[i-1]
		curr := keyPoints[i]
		distance := CalculateDistance(prev.Latitude, prev.Longitude, curr.Latitude, curr.Longitude)
		totalDistance += distance
	}

	query := "UPDATE tours SET distance_km = ?, updated_at = NOW() WHERE id = ?"
	_, err = s.db.Exec(query, totalDistance, tourID)
	if err != nil {
		return err
	}

	// Calculate automatic durations based on the distance
	return s.CalculateAutomaticDurations(uint(tourID), totalDistance)
}

func (s *Store) CalculateAutomaticDurations(tourID uint, distanceKm float64) error {
	// Delete existing durations for this tour
	deleteQuery := "DELETE FROM tour_durations WHERE tour_id = ?"
	_, err := s.db.Exec(deleteQuery, tourID)
	if err != nil {
		return err
	}

	// Calculate durations based on distance
	if distanceKm < 2.0 {
		// Less than 2km: walking only (5 km/h)
		walkingDuration := int(math.Ceil((distanceKm / 5.0) * 60)) // minutes
		insertQuery := "INSERT INTO tour_durations (tour_id, transport_type, duration_minutes, created_at) VALUES (?, ?, ?, NOW())"
		_, err = s.db.Exec(insertQuery, tourID, "WALKING", walkingDuration)
		if err != nil {
			return err
		}
	} else if distanceKm <= 5.0 {
		// 2-5km: walking and cycling
		walkingDuration := int(math.Ceil((distanceKm / 5.0) * 60))  // 5 km/h
		cyclingDuration := int(math.Ceil((distanceKm / 15.0) * 60)) // 15 km/h

		insertQuery := "INSERT INTO tour_durations (tour_id, transport_type, duration_minutes, created_at) VALUES (?, ?, ?, NOW())"

		// Insert walking duration
		_, err = s.db.Exec(insertQuery, tourID, "WALKING", walkingDuration)
		if err != nil {
			return err
		}

		// Insert cycling duration
		_, err = s.db.Exec(insertQuery, tourID, "CYCLING", cyclingDuration)
		if err != nil {
			return err
		}
	} else {
		// More than 5km: walking, cycling, and bus
		walkingDuration := int(math.Ceil((distanceKm / 5.0) * 60))  // 5 km/h
		cyclingDuration := int(math.Ceil((distanceKm / 15.0) * 60)) // 15 km/h
		busDuration := int(math.Ceil((distanceKm / 30.0) * 60))     // 30 km/h

		insertQuery := "INSERT INTO tour_durations (tour_id, transport_type, duration_minutes, created_at) VALUES (?, ?, ?, NOW())"

		// Insert walking duration
		_, err = s.db.Exec(insertQuery, tourID, "WALKING", walkingDuration)
		if err != nil {
			return err
		}

		// Insert cycling duration
		_, err = s.db.Exec(insertQuery, tourID, "CYCLING", cyclingDuration)
		if err != nil {
			return err
		}

		// Insert bus duration
		_, err = s.db.Exec(insertQuery, tourID, "BUS", busDuration)
		if err != nil {
			return err
		}
	}

	return nil
}
