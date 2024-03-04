package storage

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

func (s *RedisStorage) Increment(key string, expiration time.Duration) (int, error) {
	ctx := s.client.Context()
	result, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	err = s.client.Set(ctx, key, result, expiration).Err()
	if err != nil {
		return 0, err
	}

	return int(result), nil
}

func (s *RedisStorage) Reset(key string) error {
	ctx := s.client.Context()
	return s.client.Del(ctx, key).Err()
}
