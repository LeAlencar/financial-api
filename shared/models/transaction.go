package models

import "time"

type TransactionType string

const (
	Buy  TransactionType = "BUY"
	Sell TransactionType = "SELL"
)

type Transaction struct {
	ID           string          `json:"id" bson:"_id,omitempty"`
	UserID       string          `json:"user_id" bson:"user_id"`
	Type         TransactionType `json:"type" bson:"type"`
	CurrencyPair string          `json:"currency_pair" bson:"currency_pair"`
	Amount       float64         `json:"amount" bson:"amount"`
	ExchangeRate float64         `json:"exchange_rate" bson:"exchange_rate"`
	TotalValue   float64         `json:"total_value" bson:"total_value"`
	Status       string          `json:"status" bson:"status"`
	Timestamp    time.Time       `json:"timestamp" bson:"timestamp"`
	QuotationID  string          `json:"quotation_id" bson:"quotation_id"`
}
