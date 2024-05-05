package rabbit

import (
	"fmt"
	"os"
)

type SettingsSchema struct {
	APP_RABBIT_HOST                    string
	APP_RABBIT_PUBLISHER_EXCHANGE_NAME string
	APP_RABBIT_USERS_EXCHANGE_NAME     string
	APP_RABBIT_CONSUMER_QUEUE_NAME     string
}

func InitSettings() SettingsSchema {
	host := os.Getenv("APP_RABBIT_HOST")
	if host == "" {
		panic(fmt.Errorf("you need to specify `APP_RABBIT_HOST` environment variable"))
	}

	publisherExchangeName := os.Getenv("APP_RABBIT_PUBLISHER_EXCHANGE_NAME")
	if publisherExchangeName == "" {
		panic(fmt.Errorf("you need to specify `APP_RABBIT_PUBLISHER_EXCHANGE_NAME` environment variable"))
	}
	usersExchange := os.Getenv("APP_RABBIT_USERS_EXCHANGE_NAME")
	if usersExchange == "" {
		panic(fmt.Errorf("you need to specify `APP_RABBIT_USERS_EXCHANGE_NAME` environment variable"))
	}

	consumerQueue := os.Getenv("APP_RABBIT_CONSUMER_QUEUE_NAME")
	if consumerQueue == "" {
		panic(fmt.Errorf("you need to specify `APP_RABBIT_CONSUMER_QUEUE_NAME` environment variable"))
	}

	return SettingsSchema{
		APP_RABBIT_HOST:                    host,
		APP_RABBIT_PUBLISHER_EXCHANGE_NAME: publisherExchangeName,
		APP_RABBIT_USERS_EXCHANGE_NAME:     usersExchange,
		APP_RABBIT_CONSUMER_QUEUE_NAME:     consumerQueue,
	}
}

var Settings SettingsSchema = InitSettings()
