package kafka

import "auth/internal/domain/event"

func topicFor(e any) string {
	switch e.(type) {
	case event.UserRegistered:
		return "user.registered"
	default:
		return "unknown"
	}
}
