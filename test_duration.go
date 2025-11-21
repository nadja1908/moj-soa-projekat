package main

import (
	"database/sql"
	"log"

	"tour-service/internal/store"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Konektuj se na bazu
	dsn := "user:password@tcp(localhost:3309)/tour_db?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test konekcije
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Kreiraj store
	s := store.New(db)

	// Pozovi UpdateTourDistance za tour ID 7
	log.Println("Calling UpdateTourDistance for tour ID 7...")
	err = s.UpdateTourDistance(7)
	if err != nil {
		log.Fatal("Failed to update tour distance:", err)
	}

	log.Println("UpdateTourDistance completed successfully!")
}
