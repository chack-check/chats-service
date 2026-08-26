package main

import (
	"chats-service/internal/infrastructure/api"
	grpcservice "chats-service/internal/infrastructure/grpc_service"
	"chats-service/internal/infrastructure/rabbit"
)

func main() {
	go grpcservice.RunGrpcServer()
	rabbit.StartConsumer("chats-service")
	api.RunApi()
}
