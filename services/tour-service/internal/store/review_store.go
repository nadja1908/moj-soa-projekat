package store

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"tour-service/internal/model"
)

// ReviewStore definiše interfejs za rad sa podacima recenzija.
type ReviewStore interface {
	CreateReview(review *model.Review) error
	GetReviewsByTourID(tourID int64) ([]model.Review, error)
}

// reviewStore implementira ReviewStore koristeći *sql.DB.
type reviewStore struct {
	db *sql.DB
}

// NewReviewStore kreira novu instancu ReviewStore i tabelu reviews ako ne postoji.
func NewReviewStore(db *sql.DB) ReviewStore {
	query := `
		CREATE TABLE IF NOT EXISTS reviews (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			tour_id BIGINT NOT NULL,
			tourist_id BIGINT NOT NULL,
			rating INT NOT NULL,
			comment TEXT,
			date_visited DATE NOT NULL,
			image_urls JSON,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_tour_id (tour_id),
			INDEX idx_tourist_tour (tourist_id, tour_id)
		)
	`

	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Failed to create reviews table: %v", err)
	}

	log.Println("Database migrated for Review model (reviews table).")
	return &reviewStore{db: db}
}

// CreateReview čuva novu recenziju u bazi podataka.
func (s *reviewStore) CreateReview(review *model.Review) error {
	// ImageURLs je datatypes.JSON ([]byte) u modelu, čuvamo kao string (JSON u koloni)
	imageJSON := string(review.ImageURLs)

	// MySQL JSON kolona NE SME da dobije prazan string "" → mora biti validan JSON (npr. "[]").
	if strings.TrimSpace(imageJSON) == "" {
		imageJSON = "[]"
	}

	res, err := s.db.Exec(`
		INSERT INTO reviews (tour_id, tourist_id, rating, comment, date_visited, image_urls)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		review.TourID,
		review.TouristID,
		review.Rating,
		review.Comment,
		review.DateVisited,
		imageJSON,
	)
	if err != nil {
		log.Printf("ERROR inserting review: %v", err)
		return err
	}

	// Postavi ID nazad u struct da ga handler može vratiti klijentu
	if id, err := res.LastInsertId(); err == nil && id > 0 {
		review.ID = uint(id)
	}

	return nil
}

// pomoćna funkcija za parsiranje datuma u više formata
func parseFlexibleDate(label, value string) (time.Time, bool) {
	layouts := []string{
		"2006-01-02",          // npr. 2025-11-20 (DATE)
		"2006-01-02 15:04:05", // npr. 2025-11-24 16:55:51 (TIMESTAMP/DATETIME)
		time.RFC3339,          // npr. 2025-11-20T00:00:00Z (stari zapisi kao string)
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, true
		}
	}

	log.Printf("WARN: cannot parse %s '%s' with any known layout", label, value)
	return time.Time{}, false
}

// GetReviewsByTourID dohvaća sve recenzije za datu turu.
func (s *reviewStore) GetReviewsByTourID(tourID int64) ([]model.Review, error) {
	rows, err := s.db.Query(`
		SELECT id, tour_id, tourist_id, rating, comment, date_visited, created_at, image_urls
		FROM reviews
		WHERE tour_id = ?
		ORDER BY id DESC
	`, tourID)
	if err != nil {
		log.Printf("ERROR querying reviews: %v", err)
		return nil, err
	}
	defer rows.Close()

	var reviews []model.Review

	for rows.Next() {
		var (
			id          int64
			dbTourID    int64
			dbTouristID int64
			rating      int
			comment     sql.NullString

			// Datume čitamo kao string, pa ih fleksibilno parsiramo
			dateVisitedStr sql.NullString
			createdAtStr   sql.NullString

			imageJSON sql.NullString
		)

		if err := rows.Scan(
			&id,
			&dbTourID,
			&dbTouristID,
			&rating,
			&comment,
			&dateVisitedStr,
			&createdAtStr,
			&imageJSON,
		); err != nil {
			log.Printf("ERROR scanning review row: %v", err)
			return nil, err
		}

		var r model.Review

		// gorm.Model.ID je uint → ručno mapiramo
		if id > 0 {
			r.ID = uint(id)
		}

		r.TourID = dbTourID
		r.TouristID = dbTouristID
		r.Rating = rating
		if comment.Valid {
			r.Comment = comment.String
		}

		// date_visited
		if dateVisitedStr.Valid && dateVisitedStr.String != "" {
			if t, ok := parseFlexibleDate("date_visited", dateVisitedStr.String); ok {
				r.DateVisited = t
			}
		}

		// created_at
		if createdAtStr.Valid && createdAtStr.String != "" {
			if t, ok := parseFlexibleDate("created_at", createdAtStr.String); ok {
				r.CreatedAt = t
			}
		}

		// image_urls → uvek validan JSON niz (bar "[]")
		if imageJSON.Valid && strings.TrimSpace(imageJSON.String) != "" {
			r.ImageURLs = []byte(imageJSON.String)
		} else {
			r.ImageURLs = []byte("[]")
		}

		reviews = append(reviews, r)
	}

	if err := rows.Err(); err != nil {
		log.Printf("ERROR iterating review rows: %v", err)
		return nil, err
	}

	// Ako nema recenzija → vraćamo prazan slice, ne nil
	if reviews == nil {
		reviews = []model.Review{}
	}

	return reviews, nil
}
