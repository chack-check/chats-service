package main

import (
	"chats-service/internal/infrastructure/rabbit"
	"chats-service/internal/infrastructure/rabbit/publishers"
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := rabbit.StartConsumer(
		ctx,
		"chats-service",
		publishers.NewMessageEventsPublisher(*rabbit.EventsRabbitConnection),
	)
	if err != nil {
		log.Fatal(err)
	}
}
