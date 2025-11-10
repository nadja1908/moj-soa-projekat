package store

import (
	"database/sql"
	"fmt"

	"stakeholders-service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser kreira novog korisnika u bazi
func (s *Store) CreateUser(user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := "INSERT INTO users (username, password, email, role, is_active) VALUES (?, ?, ?, ?, ?)"
	result, err := s.db.Exec(query, user.Username, string(hashedPassword), user.Email, user.Role, true)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	user.ID = id
	user.IsActive = true
	return nil
}

// GetUserByUsername pronalazi korisnika po korisničkom imenu
func (s *Store) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	query := "SELECT id, username, password, email, role, is_active FROM users WHERE username = ?"
	
	row := s.db.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.IsActive)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Korisnik ne postoji
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	
	return user, nil
}

// GetUserByEmail pronalazi korisnika po email adresi
func (s *Store) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := "SELECT id, username, password, email, role, is_active FROM users WHERE email = ?"
	
	row := s.db.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.IsActive)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Korisnik ne postoji
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return user, nil
}

// GetUserByID pronalazi korisnika po ID-u
func (s *Store) GetUserByID(id int64) (*model.User, error) {
	user := &model.User{}
	query := "SELECT id, username, email, role, is_active FROM users WHERE id = ?"
	
	row := s.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.IsActive)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	return user, nil
}

// GetAllUsers vraća sve korisnike (za administratore)
func (s *Store) GetAllUsers() ([]model.User, error) {
	query := "SELECT id, username, email, role, is_active FROM users ORDER BY id"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return users, nil
}

// BlockUser blokira korisnikov nalog
func (s *Store) BlockUser(userID int64) error {
	query := "UPDATE users SET is_active = FALSE WHERE id = ?"
	_, err := s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}
	return nil
}

// UnblockUser odblokira korisnikov nalog
func (s *Store) UnblockUser(userID int64) error {
	query := "UPDATE users SET is_active = TRUE WHERE id = ?"
	_, err := s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to unblock user: %w", err)
	}
	return nil
}