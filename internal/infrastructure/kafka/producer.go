package kafka

import (
	"context"
	"time"

	"github.com/Ali127Dev/xerr"
	"github.com/IBM/sarama"
	"github.com/avast/retry-go"
	"go.uber.org/fx"
)

type ProducerConfig struct {
	Brokers []string
}

func NewProducer(lc fx.Lifecycle, cfg ProducerConfig) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	config.Producer.Compression = sarama.CompressionZSTD
	config.Producer.Timeout = 5 * time.Second

	var producer sarama.SyncProducer

	err := retry.Do(
		func() error {
			p, err := sarama.NewSyncProducer(cfg.Brokers, config)
			if err != nil {
				return xerr.Wrap(
					err,
					xerr.CodeServiceUnavailable,
					xerr.WithMessage("failed to connect to Kafka brokers"),
					xerr.WithMeta("brokers", cfg.Brokers),
				)
			}

			producer = p
			return nil
		},
		retry.Attempts(5),
		retry.Delay(time.Second),
		retry.LastErrorOnly(true),
	)

	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return producer.Close()
		},
	})

	return producer, nil
}
