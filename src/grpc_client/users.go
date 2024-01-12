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

type IUsersGrpc interface {
	Connect()
	GetUserByToken(token string) (*pb.UserResponse, error)
	GetUserByRefreshToken(token string) (*pb.UserResponse, error)
	GetUserById(id int) (*pb.UserResponse, error)
}

type UsersGrpc struct {
	Host string
	Port int

	client pb.UsersClient
}

type MockUsersGrpc struct{}

var UsersGrpcClient IUsersGrpc = GetUsersGrpc()

func (usersGrpc *UsersGrpc) Connect() {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	dsl := fmt.Sprintf("%s:%d", usersGrpc.Host, usersGrpc.Port)

	connection, error := grpc.Dial(dsl, opts)
	if error != nil {
		log.Printf("Error when connection to users grpc: %s", error)
	}

	usersGrpc.client = pb.NewUsersClient(connection)
}

func (usersGrpc *MockUsersGrpc) Connect() {
	log.Print("Connected to users grpc")
}

func (usersGrpc *UsersGrpc) GetUserByToken(token string) (*pb.UserResponse, error) {
	user, err := usersGrpc.client.GetUserByToken(context.Background(), &pb.GetUserByTokenRequest{Token: token})

	if err != nil || user == nil {
		return user, fmt.Errorf("Error when getting user: %v", err)
	}

	return user, nil
}

func (usersGrpc *MockUsersGrpc) GetUserByToken(token string) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Id:             1,
		Username:       "testuser",
		Phone:          "testphone",
		Email:          "testemail",
		FirstName:      "testfirstname",
		LastName:       "testlastname",
		MiddleName:     "testmiddlename",
		Activity:       "testactivity",
		Status:         "teststatus",
		EmailConfirmed: true,
		PhoneConfirmed: true,
		LastSeen:       "2023-01-01T20:00:00Z",
		AvatarUrl:      "testurl",
	}, nil
}

func (UsersGrpc *UsersGrpc) GetUserByRefreshToken(token string) (*pb.UserResponse, error) {
	user, err := UsersGrpc.client.GetUserByRefreshToken(context.Background(), &pb.GetUserByTokenRequest{Token: token})

	if err != nil || user == nil {
		return user, fmt.Errorf("Error when getting user: %v", err)
	}

	return user, nil
}

func (UsersGrpc *MockUsersGrpc) GetUserByRefreshToken(token string) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Id:             1,
		Username:       "testuser",
		Phone:          "testphone",
		Email:          "testemail",
		FirstName:      "testfirstname",
		LastName:       "testlastname",
		MiddleName:     "testmiddlename",
		Activity:       "testactivity",
		Status:         "teststatus",
		EmailConfirmed: true,
		PhoneConfirmed: true,
		LastSeen:       "2023-01-01T20:00:00Z",
		AvatarUrl:      "testurl",
	}, nil
}

func (usersGrpc *UsersGrpc) GetUserById(id int) (*pb.UserResponse, error) {
	user, err := usersGrpc.client.GetUserById(context.Background(), &pb.GetUserByIdRequest{Id: int32(id)})
	if err != nil || user == nil {
		return user, fmt.Errorf("Error when getting user: %v", err)
	}

	return user, nil
}

func (usersGrpc *MockUsersGrpc) GetUserById(id int) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Id:             1,
		Username:       "testuser",
		Phone:          "testphone",
		Email:          "testemail",
		FirstName:      "testfirstname",
		LastName:       "testlastname",
		MiddleName:     "testmiddlename",
		Activity:       "testactivity",
		Status:         "teststatus",
		EmailConfirmed: true,
		PhoneConfirmed: true,
		LastSeen:       "2023-01-01T20:00:00Z",
		AvatarUrl:      "testurl",
	}, nil
}

func GetUsersGrpc() IUsersGrpc {
	if settings.Settings.APP_ENVIRONMENT == "test" {
		return &MockUsersGrpc{}
	}

	usersGrpc := UsersGrpc{
		Host: settings.Settings.USERS_GRPC_HOST,
		Port: settings.Settings.USERS_GRPC_PORT,
	}
	usersGrpc.Connect()
	return &usersGrpc
}
