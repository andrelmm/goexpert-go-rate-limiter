package main

import (
	"context"
	"github.com/andrelmm/goexpert-go-rate-limiter/limiter"
	"github.com/andrelmm/goexpert-go-rate-limiter/middleware"
	"github.com/andrelmm/goexpert-go-rate-limiter/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}

	rdb := connectToRedis()

	st := storage.NewRedisStorage(rdb)

	rl := limiter.NewLimiter(st)

	router := gin.Default()
	router.Use(middleware.RateLimiterMiddleware(rl))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OMG it works!")
	})

	port := ":8080"
	router.Run(port)
}

func connectToRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar-se ao Redis: %v", err)
		panic(err)
	}

	log.Println("Conectado ao Redis com sucesso!")

	return rdb
}
