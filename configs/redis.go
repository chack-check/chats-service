package configs

import (
	"net/url"
	"sync"

	"github.com/caarlos0/env/v11"
)

type RedisConfiguration struct {
	URL url.URL `env:"APP_REDIS_URL,required,notEmpty"`
}

var (
	redisConfig RedisConfiguration
	redisOnce   sync.Once
)

func GetRedisConfiguration() *RedisConfiguration {
	redisOnce.Do(func() {
		if err := env.Parse(&redisConfig); err != nil {
			panic(err)
		}
	})

	return &redisConfig
}
