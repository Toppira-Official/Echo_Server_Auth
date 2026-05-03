package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/fx"
)

type KafkaProducerConfig struct {
	brokers []string
}

func NewKafkaProducer(lc fx.Lifecycle, cfg KafkaProducerConfig) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	config.Producer.Compression = sarama.CompressionZSTD
	config.Producer.Timeout = 5 * time.Second

	producer, err := sarama.NewSyncProducer(cfg.brokers, config)
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
