package event

import "time"

type UserRegistered struct {
	UserID     string
	Username   string
	OccurredAt time.Time
}
