package models

import "time"

type Quotation struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	CurrencyPair  string    `json:"currency_pair" bson:"currency_pair"` // e.g., "USD/BRL"
	BuyPrice      float64   `json:"buy_price" bson:"buy_price"`
	SellPrice     float64   `json:"sell_price" bson:"sell_price"`
	Timestamp     time.Time `json:"timestamp" bson:"timestamp"`
	LastUpdatedBy string    `json:"last_updated_by" bson:"last_updated_by"`
}
