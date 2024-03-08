package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
		ctx:    context.Background(),
	}
}

func (rs *RedisStorage) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := rs.client.Get(ctx, key)
	err := cmd.Err()
	if err != nil {
		log.Printf("Error getting value from Redis: %v", err)
	}
	return cmd
}

func (rs *RedisStorage) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := rs.client.Set(ctx, key, value, expiration)
	err := cmd.Err()
	if err != nil {
		log.Printf("Error setting value in Redis: %v", err)
	}
	return cmd
}

func (rs *RedisStorage) ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	cmd := rs.client.ZAdd(ctx, key, members...)
	err := cmd.Err()
	if err != nil {
		log.Printf("Error adding members to sorted set in Redis: %v", err)
	}
	return cmd
}

func (rs *RedisStorage) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	cmd := rs.client.ZRemRangeByScore(ctx, key, min, max)
	err := cmd.Err()
	if err != nil {
		log.Printf("Error removing members from sorted set in Redis: %v", err)
	}
	return cmd
}

func (rs *RedisStorage) ZCard(ctx context.Context, key string) *redis.IntCmd {
	cmd := rs.client.ZCard(ctx, key)
	err := cmd.Err()
	if err != nil {
		log.Printf("Error getting cardinality of sorted set from Redis: %v", err)
	}
	return cmd
}
