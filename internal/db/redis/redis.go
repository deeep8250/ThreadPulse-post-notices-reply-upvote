package redisinternal

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/threadpulse/internal/config"
)

func NewRedisClient() *redis.Client {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return client
}

func RedisInit() {
	// pinging redis
	config.RedisClient = NewRedisClient()
	ctx := context.Background()
	_, err := config.RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("redis connection failed: ", err.Error())
	}

	log.Println("Redis connected successfully")

}
