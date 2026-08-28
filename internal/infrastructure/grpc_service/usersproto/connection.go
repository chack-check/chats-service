package usersproto

import (
	"fmt"

	"chats-service/configs"
	"chats-service/internal/infrastructure/grpc_service/usersproto/usersprotobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var connection *grpc.ClientConn

func UsersClientConnect() usersprotobuf.UsersClient {
	configuration := configs.GetGRPCUsersConfiguration()
	if connection == nil || connection.GetState() != connectivity.Ready {
		opts := grpc.WithTransportCredentials(insecure.NewCredentials())
		dsl := fmt.Sprintf("%s:%d", configuration.Host, configuration.Port)

		newConnection, err := grpc.NewClient(dsl, opts)
		if err != nil {
			return nil
		}

		connection = newConnection
	}

	return usersprotobuf.NewUsersClient(connection)
}
