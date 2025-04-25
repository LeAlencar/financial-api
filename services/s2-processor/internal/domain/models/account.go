package models

import (
	"time"
)

type AccountType string

const (
	Checking   AccountType = "checking"
	Savings    AccountType = "savings"
	Credit     AccountType = "credit"
	Investment AccountType = "investment"
)

type Account struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	UserID       uint          `json:"user_id" gorm:"not null"`
	Name         string        `json:"name" gorm:"not null"`
	Type         AccountType   `json:"type" gorm:"not null"`
	Balance      float64       `json:"balance" gorm:"not null;default:0"`
	Currency     string        `json:"currency" gorm:"not null;default:'BRL'"`
	Description  string        `json:"description"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:AccountID"`
}
