package datarealm

import "time"

type Event struct {
	// Event ID
	ID int64
	// Date
	Date time.Time
	// User ID
	UserID int64
	// Type is action type
	ActionType string
	// Action in JSON
	Action string
}
