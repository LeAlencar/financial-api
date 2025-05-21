package events

type UserActionType string

const (
	UserActionCreate UserActionType = "CREATE"
	UserActionUpdate UserActionType = "UPDATE"
	UserActionDelete UserActionType = "DELETE"
)

// UserEvent represents the message structure that will be exchanged between services
type UserEvent struct {
	Action UserActionType `json:"action"`
	Data   UserEventData  `json:"data"`
}

// UserEventData contains only the necessary data for the event
type UserEventData struct {
	ID       int32  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
