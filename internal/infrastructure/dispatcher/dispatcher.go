package dispatcher

import (
	"auth/internal/domain/contract"
	"context"
	"encoding/json"

	"github.com/Ali127Dev/xoutbox"
)

type Dispatcher struct {
	store         xoutbox.Store[string]
	uuidGenerator contract.UuidGenerator
	clock         contract.Clock
}

func NewDispatcher(
	store xoutbox.Store[string],
	uuidGenerator contract.UuidGenerator,
	clock contract.Clock,
) *Dispatcher {
	return &Dispatcher{store: store, uuidGenerator: uuidGenerator, clock: clock}
}

func (k *Dispatcher) Dispatch(ctx context.Context, events ...contract.Event) error {
	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		uuid, _ := k.uuidGenerator.Generate()
		err = k.store.InsertEvent(ctx, xoutbox.Event[string]{
			ID:         uuid,
			Payload:    payload,
			EventType:  topicFor(event),
			MaxRetries: 5,
			CreatedAt:  k.clock.NowUTC(),
			Status:     xoutbox.StatusPending,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
