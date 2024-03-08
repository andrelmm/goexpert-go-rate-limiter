package limiter

import (
	"context"
	"github.com/andrelmm/goexpert-go-rate-limiter/storage"
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

func (l *Limiter) CheckRateLimit(key string, isApiKey bool) bool {
	limit, err := l.getRateLimit(isApiKey)
	if err != nil {
		log.Println("Error getting rate limit:", err)
		return false
	}

	now := time.Now().UnixMilli()

	_, err = l.storage.ZAdd(l.ctx, key, &redis.Z{Score: float64(now), Member: now}).Result()
	if err != nil {
		log.Println("Error storing request in Redis:", err)
		return false
	}

	rateLimitDuration, err := time.ParseDuration(os.Getenv("RATE_LIMIT_DURATION"))
	if err != nil {
		log.Println("Error parsing rate limit duration:", err)
		return false
	}

	_, err = l.storage.ZRemRangeByScore(l.ctx, key, "-inf", strconv.FormatInt(now-int64(rateLimitDuration.Milliseconds()), 10)).Result()
	if err != nil {
		log.Println("Error removing old requests from Redis:", err)
		return false
	}

	count, err := l.storage.ZCard(l.ctx, key).Result()
	if err != nil {
		log.Println("Error getting request count from Redis:", err)
		return false
	}

	return count <= limit
}

func (l *Limiter) getRateLimit(isAPIKey bool) (int64, error) {
	if isAPIKey {
		rateLimitToken, err := strconv.ParseInt(os.Getenv("RATE_LIMIT_TOKEN"), 10, 64)

		if err != nil {
			return 0, err
		}
		return rateLimitToken, nil
	}
	rateLimitIP, err := strconv.ParseInt(os.Getenv("RATE_LIMIT_IP"), 10, 64)
	if err != nil {
		return 0, err
	}
	return rateLimitIP, nil
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
