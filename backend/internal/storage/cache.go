package storage

import (
	"context"
	"linksnap/internal/config"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitializeRedis(env *config.Env) (*redis.Client, error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),  // "localhost:6379"
		Password: os.Getenv("REDIS_PASSWORD"), // opcional
		DB:       0,
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return RedisClient, nil
}
