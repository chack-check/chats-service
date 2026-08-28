package redisdb

import (
	"chats-service/configs"
	"context"

	"github.com/redis/go-redis/v9"
)

func InitRedisConnection() *redis.Client {
	configuration := configs.GetRedisConfiguration()
	opt, err := redis.ParseURL(configuration.URL.String())
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
