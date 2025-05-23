package models

import "time"

type TransactionType string

const (
	Withdraw TransactionType = "withdraw"
	Buy      TransactionType = "buy"
)

type Transaction struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Type      TransactionType `json:"type"`
	Currency  string          `json:"currency"`
	Amount    float64         `json:"amount"`
	Status    string          `json:"status"`
	Timestamp time.Time       `json:"timestamp"`
}
