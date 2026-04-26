package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
		return fmt.Errorf("redis get: dest cannot be nil")
	}

	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return err
		}
		return fmt.Errorf("redis get key=%s: %w", key, err)
	}

	if err := json.Unmarshal(val, dest); err != nil {
		return fmt.Errorf("redis unmarshal key=%s: %w", key, err)
	}

	return nil
}

func (r *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis marshal key=%s: %w", key, err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set key=%s: %w", key, err)
	}

	return nil
}

func (r *Cache) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete key=%s: %w", key, err)
	}
	return nil
}
