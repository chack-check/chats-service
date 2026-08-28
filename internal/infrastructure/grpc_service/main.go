package grpcservice

import (
	"fmt"
	"net"

	"chats-service/configs"
	"chats-service/internal/infrastructure/grpc_service/chatsproto"
	"chats-service/internal/infrastructure/grpc_service/chatsproto/chatsprotobuf"

	"google.golang.org/grpc"
)

func RunGrpcServer() {
	configuration := configs.GetGRPCAPIConfiguration()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", configuration.Host, configuration.Port))
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chatsServer := chatsproto.ChatsServer{}
	chatsprotobuf.RegisterChatsServer(grpcServer, chatsServer)
	grpcServer.Serve(lis)
}
