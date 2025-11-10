package model

import "database/sql"

// User predstavlja korisnički nalog u sistemu.
// Funkcionalnost #1.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Lozinka se nikada ne šalje klijentu
	Email    string `json:"email"`
	Role     string `json:"role"`     // 'guide', 'tourist', 'administrator'
	IsActive bool   `json:"isActive"` // Koristi se za blokiranje korisnika (Funkcionalnost #3)
}

// Profile predstavlja sve informacije o profilu jednog korisnika.
// Funkcionalnost #4.
type Profile struct {
	ID              int64          `json:"id"`
	UserID          int64          `json:"userId"`
	FirstName       sql.NullString `json:"firstName"`       // Može biti NULL
	LastName        sql.NullString `json:"lastName"`        // Može biti NULL
	ProfileImageURL sql.NullString `json:"profileImageUrl"` // Može biti NULL
	Biography       sql.NullString `json:"biography"`       // Može biti NULL
	Motto           sql.NullString `json:"motto"`           // Može biti NULL
}

// UserInfo je skraćena verzija za prikaz na blogu
type UserInfo struct {
	ID              int64          `json:"id"`
	Username        string         `json:"username"`
	FirstName       sql.NullString `json:"firstName"`
	ProfileImageURL sql.NullString `json:"profileImageUrl"`
}