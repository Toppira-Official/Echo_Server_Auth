package redis

import (
	"context"
	"fmt"
	"time"

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
			return retry.Do(
				func() error {
					return rdb.Ping(ctx).Err()
				},
				retry.Attempts(5),
				retry.Delay(2*time.Second),
				retry.DelayType(retry.FixedDelay),
				retry.LastErrorOnly(true),
				retry.Context(ctx),
			)
		},
		OnStop: func(ctx context.Context) error {
			return rdb.Close()
		},
	})

	return rdb
}
