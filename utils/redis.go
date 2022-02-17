package utils

import (
	"github.com/go-redis/redis/v8"
)

func InitRedisClient() *redis.Client {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       6,
		})
	return rdb
}
