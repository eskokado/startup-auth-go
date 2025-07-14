package providers

import (
	"context"
	"time"
)

type BlacklistProvider interface {
	Add(ctx context.Context, token string, ttl time.Duration) error
	Exists(ctx context.Context, token string) (bool, error)
	ExistsKey(ctx context.Context, key string) (bool, error)
	SetWithKey(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	Del(ctx context.Context, keys ...string) error
}
