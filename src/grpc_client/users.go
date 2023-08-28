package grpc_client

import (
	"context"
	"fmt"
	"log"

	pb "github.com/chack-check/chats-service/protousers"
	"github.com/chack-check/chats-service/settings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UsersGrpc struct {
	Host string
	Port int

	client pb.UsersClient
}

var UsersGrpcClient *UsersGrpc = GetUsersGrpc()

func (usersGrpc *UsersGrpc) Connect() {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	dsl := fmt.Sprintf("%s:%d", usersGrpc.Host, usersGrpc.Port)

	connection, error := grpc.Dial(dsl, opts)
	if error != nil {
		log.Printf("Error when connection to users grpc: %s", error)
	}

	usersGrpc.client = pb.NewUsersClient(connection)
}

func (usersGrpc *UsersGrpc) GetUserByToken(token string) (*pb.UserResponse, error) {
	user, err := usersGrpc.client.GetUserByToken(context.Background(), &pb.GetUserByTokenRequest{Token: token})

	if err != nil || user == nil {
		return user, fmt.Errorf("Error when getting user: %v", err)
	}

	return user, nil
}

func (UsersGrpc *UsersGrpc) GetUserByRefreshToken(token string) (*pb.UserResponse, error) {
	user, err := UsersGrpc.client.GetUserByRefreshToken(context.Background(), &pb.GetUserByTokenRequest{Token: token})

	if err != nil || user == nil {
		return user, fmt.Errorf("Error when getting user: %v", err)
	}

	return user, nil
}

func (usersGrpc *UsersGrpc) GetUserById(id int) (*pb.UserResponse, error) {
	user, err := usersGrpc.client.GetUserById(context.Background(), &pb.GetUserByIdRequest{Id: int32(id)})
	if err != nil || user == nil {
		return user, fmt.Errorf("Error when getting user: %v", err)
	}

	return user, nil
}

func GetUsersGrpc() *UsersGrpc {
	usersGrpc := UsersGrpc{
		Host: settings.Settings.USERS_GRPC_HOST,
		Port: settings.Settings.USERS_GRPC_PORT,
	}
	usersGrpc.Connect()
	return &usersGrpc
}
