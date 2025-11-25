package main

import (
	"log"

	"tour-service/internal/store"
)

func main() {
	// Kreiraj store sa istim podacima kao u main.go
	store := store.NewStore("user", "password", "localhost:3309", "tour_db")
	defer store.Close()

	// Pozovi UpdateTourDistance za tour ID 7
	log.Println("Calling UpdateTourDistance for tour ID 7...")
	err := store.UpdateTourDistance(7)
	if err != nil {
		log.Fatal("Failed to update tour distance:", err)
	}

	log.Println("UpdateTourDistance completed successfully!")
}
