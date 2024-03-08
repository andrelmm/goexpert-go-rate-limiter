package limiter

import (
	"context"
	"github.com/andrelmm/goexpert-go-rate-limiter/storage"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

var rdb *redis.Client

func TestMain(m *testing.M) {
	rdb = connectToRedis()
	defer rdb.Close()
	os.Exit(m.Run())
}

func cleanupRedis(client *redis.Client) {
	ctx := context.Background()
	iter := client.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		err := client.Del(ctx, iter.Val()).Err()
		if err != nil {
			log.Printf("Error deleting key %s: %v\n", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("Error iterating over keys: %v\n ", err)
	}
}

func connectToRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
		panic(err)
	}

	log.Println("Connected to Redis successfully!")

	return rdb
}

type MockStorage struct{}

var mockData = make(map[string]interface{})

func (m *MockStorage) ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	mockData[key] = len(members)
	return redis.NewIntResult(int64(len(members)), nil)
}

func (m *MockStorage) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	count := 0
	for k, v := range mockData {
		if k == key {
			delete(mockData, k)
			count = v.(int)
		}
	}
	return redis.NewIntResult(int64(count), nil)
}

func (m *MockStorage) ZCard(ctx context.Context, key string) *redis.IntCmd {
	count := 0
	for k, v := range mockData {
		if k == key {
			count = v.(int)
		}
	}
	return redis.NewIntResult(int64(count), nil)
}

func (m *MockStorage) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	mockData[key] = value
	return redis.NewStatusResult("OK", nil)
}

func (m *MockStorage) Get(ctx context.Context, key string) *redis.StringCmd {
	value, _ := mockData[key].(bool)
	return redis.NewStringResult(strconv.FormatBool(value), nil)
}

func TestCheckRateLimit_Positive(t *testing.T) {
	cleanupRedis(rdb)
	os.Setenv("RATE_LIMIT_DURATION", "1s")
	os.Setenv("RATE_LIMIT_TOKEN", "10")
	storage := storage.NewRedisStorage(rdb)
	lim := NewLimiter(storage)

	key := "apikey:test"
	for i := 0; i < 10; i++ {
		assert.True(t, lim.CheckRateLimit(key, true), "Request should be allowed")
		time.Sleep(100 * time.Millisecond)
	}
}

func TestCheckRateLimit_Negative(t *testing.T) {
	cleanupRedis(rdb)
	os.Setenv("RATE_LIMIT_DURATION", "10s")
	os.Setenv("RATE_LIMIT_TOKEN", "10")
	storage := storage.NewRedisStorage(rdb)
	lim := NewLimiter(storage)

	key := "apikey:test"
	for i := 0; i < 11; i++ {
		if i < 10 {
			assert.True(t, lim.CheckRateLimit(key, true), "Request should be allowed")
		} else {
			assert.False(t, lim.CheckRateLimit(key, true), "Request should be blocked")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestBlock_Positive(t *testing.T) {
	cleanupRedis(rdb)
	os.Setenv("BLOCK_DURATION", "1s")
	storage := storage.NewRedisStorage(rdb)
	lim := NewLimiter(storage)

	key := uuid.New().String()
	lim.Block(key)

	assert.True(t, lim.IsBlocked(key), "Should be blocked")
}

func TestBlock_Negative(t *testing.T) {
	cleanupRedis(rdb)
	os.Setenv("BLOCK_DURATION", "1s")
	storage := storage.NewRedisStorage(rdb)
	lim := NewLimiter(storage)

	key := uuid.New().String()
	assert.False(t, lim.IsBlocked(key), "Should not be blocked")
}
