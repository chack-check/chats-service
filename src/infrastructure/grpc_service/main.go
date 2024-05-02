package grpcservice

import (
	"fmt"
	"net"

	"github.com/chack-check/chats-service/infrastructure/grpc_service/chatsproto"
	"github.com/chack-check/chats-service/infrastructure/grpc_service/chatsproto/chatsprotobuf"
	"github.com/chack-check/chats-service/infrastructure/grpc_service/settings"
	"google.golang.org/grpc"
)

func RunGrpcServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Settings.APP_GRPC_HOST, settings.Settings.APP_GRPC_PORT))
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chatsServer := chatsproto.ChatsServer{}
	chatsprotobuf.RegisterChatsServer(grpcServer, chatsServer)
	grpcServer.Serve(lis)
}
