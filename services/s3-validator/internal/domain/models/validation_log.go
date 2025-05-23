package models

import (
	"time"

	"github.com/gocql/gocql"
)

// EventType represents the type of event being logged
type EventType string

const (
	EventTypeUser        EventType = "user"
	EventTypeTransaction EventType = "transaction"
	EventTypeQuotation   EventType = "quotation"
)

// ActionType represents the action being performed
type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
	ActionBuy    ActionType = "buy"
	ActionSell   ActionType = "sell"
)

// ValidationStatus represents the validation result
type ValidationStatus string

const (
	ValidationSuccess ValidationStatus = "success"
	ValidationFailed  ValidationStatus = "failed"
	ValidationError   ValidationStatus = "error"
)

// ValidationLog represents a comprehensive log entry for any event in the system
type ValidationLog struct {
	ID          gocql.UUID        `json:"id"`
	EventType   EventType         `json:"event_type"`
	Action      ActionType        `json:"action"`
	EventID     string            `json:"event_id"`
	UserID      string            `json:"user_id,omitempty"`
	Status      ValidationStatus  `json:"status"`
	Message     string            `json:"message"`
	Details     map[string]string `json:"details,omitempty"`
	RawPayload  string            `json:"raw_payload"`
	ProcessedAt time.Time         `json:"processed_at"`
	Source      string            `json:"source"` // Which service sent the event
}

// MessageLog represents a simplified log for message tracking
type MessageLog struct {
	ID          gocql.UUID `json:"id"`
	QueueName   string     `json:"queue_name"`
	EventType   EventType  `json:"event_type"`
	Action      ActionType `json:"action"`
	Status      string     `json:"status"`
	RawPayload  string     `json:"raw_payload"`
	ProcessedAt time.Time  `json:"processed_at"`
	ErrorMsg    string     `json:"error_msg,omitempty"`
}

// ValidationRule represents a rule for validating events
type ValidationRule struct {
	EventType EventType  `json:"event_type"`
	Action    ActionType `json:"action"`
	Field     string     `json:"field"`
	Rule      string     `json:"rule"`
	Required  bool       `json:"required"`
}
