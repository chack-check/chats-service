package redisdb

import (
	"github.com/chack-check/chats-service/settings"
	"github.com/redis/go-redis/v9"
)

func InitRedisConnection() *redis.Client {
	opt, err := redis.ParseURL(settings.Settings.APP_REDIS_URL)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	return client
}

var RedisConnection = InitRedisConnection()
