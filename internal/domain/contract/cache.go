package contract

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (value any, err error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) (err error)
	Delete(ctx context.Context, key string) (err error)
}
