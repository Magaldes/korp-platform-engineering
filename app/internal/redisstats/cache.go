package redisstats

import (
	"context"
	"time"

	"github.com/projeto-korp/app/internal/audit"
	"github.com/redis/go-redis/v9"
)

type Cache struct{ client *redis.Client }

func New(address string) *Cache {
	return &Cache{client: redis.NewClient(&redis.Options{Addr: address})}
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	payload, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, audit.ErrCacheMiss
	}
	return payload, err
}

func (c *Cache) Set(ctx context.Context, key string, payload []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, payload, ttl).Err()
}

func (c *Cache) Close() error { return c.client.Close() }

var _ audit.StatisticsCache = (*Cache)(nil)
