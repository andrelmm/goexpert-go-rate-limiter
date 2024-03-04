package limiter

import (
	"context"
	"github.com/andrelmm/goexpert-go-rate-limiter/pkg/storage"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"time"
)

type Limiter struct {
	storage storage.Storage
	ctx     context.Context
}

func NewLimiter(storage storage.Storage) *Limiter {
	return &Limiter{
		storage: storage,
		ctx:     context.Background(),
	}
}

func (l *Limiter) CheckRateLimit(key string) bool {
	limit, err := l.storage.Get(l.ctx, key).Int64()
	if err != nil {
		limit = 0
	}

	if limit == 0 {
		return true
	}

	now := time.Now().Unix()

	_, err = l.storage.ZAdd(l.ctx, key, &redis.Z{Score: float64(now), Member: now}).Result()
	if err != nil {
		log.Println("Error storing request in Redis:", err)
		return false
	}

	rateLimitDuration, err := time.ParseDuration(os.Getenv("RATE_LIMIT_DURATION"))
	l.storage.ZRemRangeByScore(l.ctx, key, "-inf", strconv.FormatInt(now-int64(rateLimitDuration), 10))

	count, err := l.storage.ZCard(l.ctx, key).Result()
	if err != nil {
		log.Println("Error getting request count from Redis:", err)
		return false
	}

	return count <= limit
}

func (l *Limiter) Block(key string) {
	duration, err := time.ParseDuration(os.Getenv("BLOCK_DURATION"))
	if err != nil {
		log.Println("Error parsing block duration:", err)
		return
	}
	err = l.storage.Set(l.ctx, key+":blocked", true, duration).Err()
	if err != nil {
		log.Println("Error storing block in Redis:", err)
		return
	}
}

func (l *Limiter) IsBlocked(key string) bool {
	blocked, err := l.storage.Get(l.ctx, key+":blocked").Bool()
	if err != nil {
		return false
	}
	return blocked
}
