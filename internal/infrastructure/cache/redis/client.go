package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Ali127Dev/xerr"
	"github.com/avast/retry-go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type ClientConfig struct {
	Host     string
	Password string
	Port     int
	DB       int
}

func NewClient(lc fx.Lifecycle, cfg ClientConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		DialTimeout:     3 * time.Second,
		ReadTimeout:     1 * time.Second,
		WriteTimeout:    1 * time.Second,
		PoolSize:        100,
		MinIdleConns:    20,
		PoolTimeout:     1 * time.Second,
		MaxRetries:      2,
		MinRetryBackoff: 100 * time.Millisecond,
		MaxRetryBackoff: 500 * time.Millisecond,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := retry.Do(
				func() error {
					if err := rdb.Ping(ctx).Err(); err != nil {
						return xerr.Wrap(
							err,
							xerr.CodeServiceUnavailable,
							xerr.WithMessage("redis is not reachable"),
						)
					}
					return nil
				},
				retry.Attempts(5),
				retry.Delay(2*time.Second),
				retry.DelayType(retry.FixedDelay),
				retry.LastErrorOnly(true),
				retry.Context(ctx),
			)

			if err != nil {
				return xerr.Wrap(
					err,
					xerr.CodeServiceUnavailable,
					xerr.WithMessage("redis connection retries exhausted"),
				)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := rdb.Close(); err != nil {
				return xerr.Wrap(
					err,
					xerr.CodeInternalError,
					xerr.WithMessage("failed to close redis client"),
				)
			}
			return nil
		},
	})

	return rdb
}
