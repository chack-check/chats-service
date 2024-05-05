package usersproto

import (
	"fmt"

	"github.com/chack-check/chats-service/infrastructure/grpc_service/usersproto/usersprotobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var connection *grpc.ClientConn

func UsersClientConnect() usersprotobuf.UsersClient {
	if connection == nil || connection.GetState() != connectivity.Ready {
		opts := grpc.WithTransportCredentials(insecure.NewCredentials())
		dsl := fmt.Sprintf("%s:%d", Settings.APP_USERS_GRPC_HOST, Settings.APP_USERS_GRPC_PORT)

		newConnection, err := grpc.Dial(dsl, opts)
		if err != nil {
			return nil
		}

		connection = newConnection
	}

	return usersprotobuf.NewUsersClient(connection)
}
