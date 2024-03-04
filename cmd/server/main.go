package server

import (
	"github.com/andrelmm/goexpert-go-rate-limiter/pkg/limiter"
	"github.com/andrelmm/goexpert-go-rate-limiter/pkg/middleware"
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

	rl := limiter.NewLimiter(rdb)

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
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}
