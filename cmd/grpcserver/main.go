package main

import (
	grpcservice "chats-service/internal/infrastructure/grpc_service"
)

func main() {
	grpcservice.RunGrpcServer()
}
