package outbox

import (
	"context"

	"github.com/Ali127Dev/xoutbox"
	"go.uber.org/fx"
)

func Start(
	lc fx.Lifecycle,
	store xoutbox.Store[string],
	publisher xoutbox.Publisher[string],
) error {
	worker := xoutbox.NewWorker(
		store,
		publisher,
		xoutbox.DefaultWorkerConfig(),
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return worker.Start(ctx)
		},
	})

	return nil
}
