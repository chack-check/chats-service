package main

import (
	"github.com/chack-check/chats-service/infrastructure/api"
	grpcservice "github.com/chack-check/chats-service/infrastructure/grpc_service"
	"github.com/chack-check/chats-service/infrastructure/rabbit"
)

func main() {
	go grpcservice.RunGrpcServer()
	rabbit.StartConsumer("chats-service")
	api.RunApi()
}
