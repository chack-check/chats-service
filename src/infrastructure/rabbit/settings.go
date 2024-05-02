package rabbit

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type SettingsSchema struct {
	APP_RABBIT_USER                    string
	APP_RABBIT_PASSWORD                string
	APP_RABBIT_HOST                    string
	APP_RABBIT_PORT                    int
	APP_RABBIT_PUBLISHER_EXCHANGE_NAME string
	APP_RABBIT_USERS_EXCHANGE_NAME     string
	APP_RABBIT_CONSUMER_QUEUE_NAME     string
}

func InitSettings() SettingsSchema {
	user := os.Getenv("APP_RABBIT_USER")
	if user == "" {
		panic(fmt.Errorf("you need to specify `APP_RABBIT_USER` environment variable"))
	}

	password := os.Getenv("APP_RABBIT_PASSWORD")
	if password == "" {
		panic(fmt.Errorf("you need to specify `APP_RABBIT_PASSWORD` environment variable"))
	}

	host := os.Getenv("APP_RABBIT_HOST")
	if host == "" {
		panic(fmt.Errorf("you need to specify `APP_RABBIT_HOST` environment variable"))
	}

	port := os.Getenv("APP_RABBIT_PORT")
	if port == "" {
		panic(fmt.Errorf("you need to specify `APP_RABBIT_PORT` environment variable"))
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(errors.Join(fmt.Errorf("specify the correct number for rabbit port"), err))
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
		APP_RABBIT_USER:                    user,
		APP_RABBIT_PASSWORD:                password,
		APP_RABBIT_HOST:                    host,
		APP_RABBIT_PORT:                    portInt,
		APP_RABBIT_PUBLISHER_EXCHANGE_NAME: publisherExchangeName,
		APP_RABBIT_USERS_EXCHANGE_NAME:     usersExchange,
		APP_RABBIT_CONSUMER_QUEUE_NAME:     consumerQueue,
	}
}

var Settings SettingsSchema = InitSettings()
