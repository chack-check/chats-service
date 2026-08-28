package configs

import (
	"sync"

	"github.com/caarlos0/env/v11"
)

type RabbitConfiguration struct {
	Host                    string `env:"APP_RABBIT_HOST,required,notEmpty"`
	PublisherExchangeName   string `env:"APP_RABBIT_PUBLISHER_EXCHANGE_NAME,required,notEmpty"`
	UsersExchangeName       string `env:"APP_RABBIT_USERS_EXCHANGE_NAME,required,notEmpty"`
	ConsumerQueueName       string `env:"APP_RABBIT_CONSUMER_QUEUE_NAME,required,notEmpty"`
	RecognitionQueueName    string `env:"APP_RABBIT_RECOGNITION_QUEUE_NAME,required,notEmpty"`
	RecognitionExchangeName string `env:"APP_RABBIT_RECOGNITION_EXCHANGE_NAME,required,notEmpty"`
}

var (
	rabbitConfig     RabbitConfiguration
	rabbitConfigOnce sync.Once
)

func GetRabbitConfiguration() *RabbitConfiguration {
	rabbitConfigOnce.Do(func() {
		if err := env.Parse(&rabbitConfig); err != nil {
			panic(err)
		}
	})

	return &rabbitConfig
}
