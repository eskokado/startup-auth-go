package providers

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCmdable define a interface mínima necessária para o blacklist
type RedisCmdable interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
}

type RedisBlacklist struct {
	client RedisCmdable
}

func NewRedisBlacklist(client RedisCmdable) *RedisBlacklist {
	return &RedisBlacklist{client: client}
}

func (r *RedisBlacklist) Add(
	ctx context.Context,
	token string,
	ttl time.Duration,
) error {
	value := ""
	if ttl > 0 {
		value = token
	}
	return r.client.Set(ctx, "blacklist:"+token, value, ttl).Err()
}

func (r *RedisBlacklist) Exists(ctx context.Context, token string) (bool, error) {
	exists, err := r.client.Exists(ctx, "blacklist:"+token).Result()
	return exists > 0, err
}
