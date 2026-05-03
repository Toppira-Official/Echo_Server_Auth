package contract

import "context"

type Event any

type EventDispatcher interface {
	Dispatch(ctx context.Context, events ...Event) error
}
