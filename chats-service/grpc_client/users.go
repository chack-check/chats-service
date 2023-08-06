package grpc_client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	pb "github.com/chack-check/chats-service/protousers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UsersGrpc struct {
	Host string
	Port int

	client pb.UsersClient
}

func (usersGrpc *UsersGrpc) Connect() {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	dsl := fmt.Sprintf("%s:%d", usersGrpc.Host, usersGrpc.Port)

	connection, error := grpc.Dial(dsl, opts)
	if error != nil {
		log.Fatalf("Error when connection to users grpc: %s", error)
	}

	usersGrpc.client = pb.NewUsersClient(connection)
}

func (usersGrpc *UsersGrpc) GetUserByToken(token string) (*pb.UserResponse, error) {
	user, err := usersGrpc.client.GetUserByToken(context.Background(), &pb.GetUserByTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (UsersGrpc *UsersGrpc) GetUserById(id int) (*pb.UserResponse, error) {
	user, err := UsersGrpc.client.GetUserById(context.Background(), &pb.GetUserByIdRequest{Id: int32(id)})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUsersGrpc() *UsersGrpc {
	host := os.Getenv("USERS_GRPC_HOST")
	port, err := strconv.Atoi(os.Getenv("USERS_GRPC_PORT"))
	if err != nil {
		log.Fatalf("Error when parsing USERS_GRPC_PORT env var: %s", err)
	}

	usersGrpc := UsersGrpc{Host: host, Port: port}
	usersGrpc.Connect()
	return &usersGrpc
}
