package main

import (
	"chats-service/internal/infrastructure/api"
	grpcservice "chats-service/internal/infrastructure/grpc_service"
	"chats-service/internal/infrastructure/rabbit"
	"chats-service/internal/infrastructure/rabbit/publishers"
)

func main() {
	// TODO: Make 3 cmds and 3 differrent build stages
	go grpcservice.RunGrpcServer()
	rabbit.StartConsumer(
		"chats-service",
		publishers.NewMessageEventsPublisher(*rabbit.EventsRabbitConnection),
	)
	api.RunApi()
}
