package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Storage interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd
	ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd
	ZCard(ctx context.Context, key string) *redis.IntCmd
}
