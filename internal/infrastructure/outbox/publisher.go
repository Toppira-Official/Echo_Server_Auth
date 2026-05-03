package outbox

import (
	"time"

	"github.com/Ali127Dev/xerr"
	"github.com/Ali127Dev/xoutbox"
	"github.com/Ali127Dev/xoutbox/kafka"
	"github.com/IBM/sarama"
	"go.uber.org/fx"
)

type PublisherConfig struct {
	Brokers []string
}

func NewPublisher(lc fx.Lifecycle, cfg PublisherConfig) (xoutbox.Publisher[string], error) {
	publisher, err := kafka.NewPublisher[string](kafka.Config{
		Brokers:      cfg.Brokers,
		RequiredAcks: sarama.WaitForAll,
		BatchSize:    1,
		BatchTimeout: 1 * time.Second,
		Topic:        "user.registered",
	})
	if err != nil {
		return nil, xerr.Wrap(
			err,
			xerr.CodeServiceUnavailable,
			xerr.WithMessage("failed to open kafka connection"),
		)
	}

	lc.Append(fx.StopHook(publisher.Close))

	return publisher, nil
}
