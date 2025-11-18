package store

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbUser, dbPass, dbHost, dbName string) *Store {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	store := &Store{db: db}
	store.createTables()
	return store
}

func (s *Store) createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tours (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			difficulty ENUM('easy', 'medium', 'hard') NOT NULL,
			tags JSON,
			status ENUM('draft', 'published', 'archived') NOT NULL DEFAULT 'draft',
			price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
			author_id BIGINT NOT NULL,
			distance_km DECIMAL(8,2) NOT NULL DEFAULT 0.00,
			published_at TIMESTAMP NULL,
			archived_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_author_id (author_id),
			INDEX idx_status (status)
		)`,

		`CREATE TABLE IF NOT EXISTS key_points (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			tour_id BIGINT NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			image_url VARCHAR(500),
			latitude DECIMAL(10,8) NOT NULL,
			longitude DECIMAL(11,8) NOT NULL,
			order_index INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (tour_id) REFERENCES tours(id) ON DELETE CASCADE,
			INDEX idx_tour_id (tour_id),
			INDEX idx_order (tour_id, order_index)
		)`,

		`CREATE TABLE IF NOT EXISTS tour_durations (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			tour_id BIGINT NOT NULL,
			transport_type ENUM('walk', 'bicycle', 'car') NOT NULL,
			duration_minutes INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (tour_id) REFERENCES tours(id) ON DELETE CASCADE,
			UNIQUE KEY unique_tour_transport (tour_id, transport_type),
			INDEX idx_tour_id (tour_id)
		)`,

		`CREATE TABLE IF NOT EXISTS tour_executions (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			tour_id BIGINT NOT NULL,
			tourist_id BIGINT NOT NULL,
			current_latitude DECIMAL(10,8),
			current_longitude DECIMAL(11,8),
			started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP NULL,
			FOREIGN KEY (tour_id) REFERENCES tours(id) ON DELETE CASCADE,
			INDEX idx_tour_id (tour_id),
			INDEX idx_tourist_id (tourist_id)
		)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}
}

func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}
