package models

import "time"

type UserActionType string

const (
	UserActionCreate UserActionType = "CREATE"
	UserActionUpdate UserActionType = "UPDATE"
	UserActionDelete UserActionType = "DELETE"
)

type UserMessage struct {
	Action UserActionType `json:"action"`
	User   *User          `json:"user"`
}

type User struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // "-" ensures password is never sent in JSON responses
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
