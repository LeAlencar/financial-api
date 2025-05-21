package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents the complete user model for write operations
type User struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Password is never exposed in JSON
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository defines the interface for user write operations
type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	Delete(id int32) error
	GetByID(id int32) (*User, error)        // Needed for validation
	GetByEmail(email string) (*User, error) // Needed for validation
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
