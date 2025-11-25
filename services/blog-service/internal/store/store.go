package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbUser, dbPass, dbHost, dbName string) (*Store, error) {

	// Dobavi port iz env varijable
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		// Ako nije postavljen, koristi default (za Docker)
		dbPort = "3306"
	}

	// Format DSN stringa
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connected to MySQL at %s:%s\n", dbHost, dbPort)
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
