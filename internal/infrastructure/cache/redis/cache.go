package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Ali127Dev/xerr"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (r *Cache) Get(ctx context.Context, key string, dest any) (err error) {
	if dest == nil {
		return xerr.New(
			xerr.CodeInternalError,
			xerr.WithMessage("redis get destination cannot be nil"),
		)
	}

	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return xerr.New(
				xerr.CodeNotFound,
				xerr.WithMessage("cache key not found"),
				xerr.WithMeta("key", key),
			)
		}

		return xerr.Wrap(
			err,
			xerr.CodeServiceUnavailable,
			xerr.WithMessage("failed to get value from cache"),
			xerr.WithMeta("key", key),
		)
	}

	if err := json.Unmarshal(val, dest); err != nil {
		return xerr.Wrap(
			err,
			xerr.CodeInternalError,
			xerr.WithMessage("failed to decode cached value"),
			xerr.WithMeta("key", key),
		)
	}

	return nil
}

func (r *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return xerr.Wrap(
			err,
			xerr.CodeInternalError,
			xerr.WithMessage("failed to encode cache value"),
			xerr.WithMeta("key", key),
		)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return xerr.Wrap(
			err,
			xerr.CodeServiceUnavailable,
			xerr.WithMessage("failed to set cache value"),
			xerr.WithMeta("key", key),
		)
	}

	return nil
}

func (r *Cache) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return xerr.Wrap(
			err,
			xerr.CodeServiceUnavailable,
			xerr.WithMessage("failed to delete cache key"),
			xerr.WithMeta("key", key),
		)
	}

	return nil
}
