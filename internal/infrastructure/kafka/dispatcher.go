package kafka

import (
	"auth/internal/domain/contract"
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
)

type Dispatcher struct {
	producer sarama.SyncProducer
}

func NewDispatcher(producer sarama.SyncProducer) *Dispatcher {
	return &Dispatcher{producer: producer}
}

func (k *Dispatcher) Dispatch(ctx context.Context, events ...contract.Event) error {
	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		msg := &sarama.ProducerMessage{
			Topic: topicFor(event),
			Value: sarama.ByteEncoder(payload),
		}

		_, _, err = k.producer.SendMessage(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
