package store

import (
	"database/sql"
	"fmt"

	"tour-service/internal/model"
)

// CreateKeyPoint kreira novu ključnu tačku
func (s *Store) CreateKeyPoint(req *model.CreateKeyPointRequest) (*model.KeyPoint, error) {
	// Dobij sledeći order_index za ovu turu
	orderIndex, err := s.getNextOrderIndex(req.TourID)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO key_points (tour_id, name, description, latitude, longitude, order_index)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(query, req.TourID, req.Name, req.Description,
		req.Latitude, req.Longitude, orderIndex)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	keyPoint, err := s.GetKeyPointByID(id)
	if err != nil {
		return nil, err
	}

	// Ažuriraj dužinu ture nakon dodavanja nove tačke
	s.UpdateTourDistance(req.TourID)

	return keyPoint, nil
}

// GetKeyPointByID vraća ključnu tačku po ID-u
func (s *Store) GetKeyPointByID(id int64) (*model.KeyPoint, error) {
	query := `
		SELECT id, tour_id, name, description, latitude, longitude,
		       order_index, created_at, updated_at
		FROM key_points WHERE id = ?
	`

	row := s.db.QueryRow(query, id)
	keyPoint := &model.KeyPoint{}

	err := row.Scan(
		&keyPoint.ID, &keyPoint.TourID, &keyPoint.Name, &keyPoint.Description,
		&keyPoint.Latitude, &keyPoint.Longitude,
		&keyPoint.OrderIndex, &keyPoint.CreatedAt, &keyPoint.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return keyPoint, nil
}

// GetKeyPoints vraća sve ključne tačke za turu, sortirane po order_index
func (s *Store) GetKeyPoints(tourID int64) ([]model.KeyPoint, error) {
	query := `
		SELECT id, tour_id, name, description, latitude, longitude,
		       order_index, created_at, updated_at
		FROM key_points 
		WHERE tour_id = ? 
		ORDER BY order_index ASC
	`

	rows, err := s.db.Query(query, tourID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keyPoints []model.KeyPoint
	for rows.Next() {
		keyPoint := model.KeyPoint{}
		err := rows.Scan(
			&keyPoint.ID, &keyPoint.TourID, &keyPoint.Name, &keyPoint.Description,
			&keyPoint.Latitude, &keyPoint.Longitude,
			&keyPoint.OrderIndex, &keyPoint.CreatedAt, &keyPoint.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		keyPoints = append(keyPoints, keyPoint)
	}

	return keyPoints, nil
}

// UpdateKeyPoint ažurira ključnu tačku
func (s *Store) UpdateKeyPoint(id int64, req *model.UpdateKeyPointRequest) (*model.KeyPoint, error) {
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
	if req.Latitude != nil {
		setParts = append(setParts, "latitude = ?")
		args = append(args, *req.Latitude)
	}
	if req.Longitude != nil {
		setParts = append(setParts, "longitude = ?")
		args = append(args, *req.Longitude)
	}

	if len(setParts) == 0 {
		return s.GetKeyPointByID(id)
	}

	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE key_points SET %s WHERE id = ?", setParts[0])
	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf("%s, %s", query, setParts[i])
	}

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	keyPoint, err := s.GetKeyPointByID(id)
	if err != nil {
		return nil, err
	}

	// Ažuriraj dužinu ture nakon izmene tačke
	s.UpdateTourDistance(keyPoint.TourID)

	return keyPoint, nil
}

// DeleteKeyPoint briše ključnu tačku
func (s *Store) DeleteKeyPoint(id int64) error {
	// Prvo dobij tour_id za ažuriranje distance
	keyPoint, err := s.GetKeyPointByID(id)
	if err != nil {
		return err
	}
	tourID := keyPoint.TourID

	query := "DELETE FROM key_points WHERE id = ?"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("key point not found")
	}

	// Ažuriraj order_index za preostale tačke
	err = s.reorderKeyPoints(tourID)
	if err != nil {
		return err
	}

	// Ažuriraj dužinu ture nakon brisanja tačke
	s.UpdateTourDistance(tourID)

	return nil
}

// ReorderKeyPoints reorganizuje redosled ključnih tačaka
func (s *Store) ReorderKeyPoints(tourID int64, keyPointIDs []int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, keyPointID := range keyPointIDs {
		query := "UPDATE key_points SET order_index = ?, updated_at = NOW() WHERE id = ? AND tour_id = ?"
		_, err := tx.Exec(query, i+1, keyPointID, tourID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	// Ažuriraj dužinu ture nakon promene redosleda
	s.UpdateTourDistance(tourID)

	return nil
}

// getNextOrderIndex vraća sledeći order_index za novu ključnu tačku
func (s *Store) getNextOrderIndex(tourID int64) (int, error) {
	var maxOrder sql.NullInt64
	query := "SELECT MAX(order_index) FROM key_points WHERE tour_id = ?"
	err := s.db.QueryRow(query, tourID).Scan(&maxOrder)
	if err != nil {
		return 0, err
	}

	if maxOrder.Valid {
		return int(maxOrder.Int64) + 1, nil
	}
	return 1, nil
}

// reorderKeyPoints reorganizuje order_index nakon brisanja
func (s *Store) reorderKeyPoints(tourID int64) error {
	keyPoints, err := s.GetKeyPoints(tourID)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, keyPoint := range keyPoints {
		query := "UPDATE key_points SET order_index = ?, updated_at = NOW() WHERE id = ?"
		_, err := tx.Exec(query, i+1, keyPoint.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
