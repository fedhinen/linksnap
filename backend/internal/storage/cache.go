package storage

import (
	"context"
	"linksnap/internal/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitializeRedis(env *config.Env) (*redis.Client, error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     env.RedisURL,      // "localhost:6379"
		Password: env.RedisPassword, // opcional
		Username: env.RedisUsername,
		DB:       0,
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return RedisClient, nil
}
