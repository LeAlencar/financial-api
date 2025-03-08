package models

import (
	"time"
)

type TransactionType string

const (
	Income   TransactionType = "income"
	Expense  TransactionType = "expense"
	Transfer TransactionType = "transfer"
)

type Transaction struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	AccountID   uint            `json:"account_id" gorm:"not null"`
	CategoryID  uint            `json:"category_id"`
	Amount      float64         `json:"amount" gorm:"not null"`
	Type        TransactionType `json:"type" gorm:"not null"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date" gorm:"not null"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
