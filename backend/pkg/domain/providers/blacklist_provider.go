package providers

import (
	"context"
	"time"
)

type BlacklistProvider interface {
	Add(ctx context.Context, token string, ttl time.Duration) error
	Exists(ctx context.Context, token string) (bool, error)
}
