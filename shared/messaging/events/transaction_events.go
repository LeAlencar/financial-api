package events

import "time"

type TransactionActionType string

const (
	TransactionActionBuy  TransactionActionType = "BUY"
	TransactionActionSell TransactionActionType = "SELL"
)

// TransactionEvent represents the message structure for transaction operations
type TransactionEvent struct {
	Action TransactionActionType `json:"action"`
	Data   TransactionEventData  `json:"data"`
}

// TransactionEventData contains the necessary data for transaction events
type TransactionEventData struct {
	UserID       string    `json:"user_id"`
	CurrencyPair string    `json:"currency_pair"`
	Amount       float64   `json:"amount"`
	QuotationID  string    `json:"quotation_id,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}
