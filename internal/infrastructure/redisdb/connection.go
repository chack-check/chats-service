package redisdb

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func InitRedisConnection() *redis.Client {
	opt, err := redis.ParseURL(Settings.APP_REDIS_URL)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic("Error connecting to redis")
	}

	return client
}

var RedisConnection = InitRedisConnection()
