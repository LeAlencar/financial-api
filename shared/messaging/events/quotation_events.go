package events

import "time"

type QuotationActionType string

const (
	QuotationActionCreate QuotationActionType = "CREATE_QUOTATION"
	QuotationActionUpdate QuotationActionType = "UPDATE_QUOTATION"
	QuotationActionDelete QuotationActionType = "DELETE_QUOTATION"
)

// UserEvent represents the message structure that will be exchanged between services
type QuotationEvent struct {
	Action QuotationActionType `json:"action"`
	Data   QuotationEventData  `json:"data"`
}

// UserEventData contains only the necessary data for the event
type QuotationEventData struct {
	CurrencyPair string    `json:"currency_pair" bson:"currency_pair"` // e.g., "USD/BRL"
	BuyPrice     float64   `json:"buy_price" bson:"buy_price"`
	SellPrice    float64   `json:"sell_price" bson:"sell_price"`
	Timestamp    time.Time `json:"timestamp" bson:"timestamp"`
}
