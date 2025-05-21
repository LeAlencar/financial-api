package models

import "time"

// User represents the read model for users in the generator service
type User struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository defines the interface for user read operations
type UserRepository interface {
	GetByID(id int32) (*User, error)
	GetByEmail(email string) (*User, error)
}
