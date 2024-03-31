package grpc_client

import (
	"context"
	"fmt"
	"log"

	pb "github.com/chack-check/chats-service/protousers"
	"github.com/chack-check/chats-service/settings"
	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IUsersGrpc interface {
	Connect()
	GetUserByToken(token string) (*pb.UserResponse, error)
	GetUserByRefreshToken(token string) (*pb.UserResponse, error)
	GetUserById(id int) (*pb.UserResponse, error)
	GetUsersByIds(ids []int) (*pb.UsersArrayResponse, error)
}

type UsersGrpc struct {
	Host string
	Port int

	client pb.UsersClient
}

type MockUsersGrpc struct{}

var UsersGrpcClient IUsersGrpc = GetUsersGrpc()

func (usersGrpc *UsersGrpc) Connect() {
	log.Printf("Connecting to users grpc host = %s port = %d", usersGrpc.Host, usersGrpc.Port)
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	dsl := fmt.Sprintf("%s:%d", usersGrpc.Host, usersGrpc.Port)

	connection, err := grpc.Dial(dsl, opts)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Error when connection to users grpc: %s", err)
		return
	}

	usersGrpc.client = pb.NewUsersClient(connection)
}

func (usersGrpc *MockUsersGrpc) Connect() {
	log.Print("Connected to users grpc")
}

func (usersGrpc *UsersGrpc) GetUserByToken(token string) (*pb.UserResponse, error) {
	log.Printf("Fetching user by token from users grpc")
	user, err := usersGrpc.client.GetUserByToken(context.Background(), &pb.GetUserByTokenRequest{Token: token})

	if err != nil || user == nil {
		if err != nil {
			sentry.CaptureException(err)
		} else {
			sentry.CaptureException(fmt.Errorf("error when getting user from users grpc"))
		}
		return user, fmt.Errorf("error when getting user: %v", err)
	}

	return user, nil
}

func (usersGrpc *MockUsersGrpc) GetUserByToken(token string) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Id:                 1,
		Username:           "testuser",
		Phone:              "testphone",
		Email:              "testemail",
		FirstName:          "testfirstname",
		LastName:           "testlastname",
		MiddleName:         "testmiddlename",
		Activity:           "testactivity",
		Status:             "teststatus",
		EmailConfirmed:     true,
		PhoneConfirmed:     true,
		LastSeen:           "2023-01-01T20:00:00Z",
		OriginalAvatarUrl:  "testurl",
		ConvertedAvatarUrl: "testurl",
	}, nil
}

func (UsersGrpc *UsersGrpc) GetUserByRefreshToken(token string) (*pb.UserResponse, error) {
	log.Printf("Fetching user by refresh token from users grpc")
	user, err := UsersGrpc.client.GetUserByRefreshToken(context.Background(), &pb.GetUserByTokenRequest{Token: token})

	if err != nil || user == nil {
		log.Printf("Error when getting user: %v", err)

		if err != nil {
			sentry.CaptureException(err)
		} else {
			sentry.CaptureException(fmt.Errorf("error when getting user by refresh token"))
		}

		return user, fmt.Errorf("error when getting user: %v", err)
	}

	return user, nil
}

func (UsersGrpc *MockUsersGrpc) GetUserByRefreshToken(token string) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Id:                 1,
		Username:           "testuser",
		Phone:              "testphone",
		Email:              "testemail",
		FirstName:          "testfirstname",
		LastName:           "testlastname",
		MiddleName:         "testmiddlename",
		Activity:           "testactivity",
		Status:             "teststatus",
		EmailConfirmed:     true,
		PhoneConfirmed:     true,
		LastSeen:           "2023-01-01T20:00:00Z",
		OriginalAvatarUrl:  "testurl",
		ConvertedAvatarUrl: "testurl",
	}, nil
}

func (usersGrpc *UsersGrpc) GetUserById(id int) (*pb.UserResponse, error) {
	log.Printf("Fetching user by id from users grpc: %d", id)
	user, err := usersGrpc.client.GetUserById(context.Background(), &pb.GetUserByIdRequest{Id: int32(id)})
	if err != nil || user == nil {
		log.Printf("Error fetching user by id from users grpc: %v", err)
		if err != nil {
			sentry.CaptureException(err)
		} else {
			sentry.CaptureException(fmt.Errorf("error fetching user by id"))
		}

		return user, fmt.Errorf("error when getting user: %v", err)
	}

	return user, nil
}

func (usersGrpc *UsersGrpc) GetUsersByIds(ids []int) (*pb.UsersArrayResponse, error) {
	log.Printf("Fetching users by ids from users grpc: %v", ids)
	var ids_int32 []int32
	for _, id := range ids {
		ids_int32 = append(ids_int32, int32(id))
	}

	users_response, err := usersGrpc.client.GetUsersByIds(context.Background(), &pb.GetUsersByIdsRequest{Ids: ids_int32})
	if err != nil {
		log.Printf("Error fetching users by ids from users grpc: %v", err)
		sentry.CaptureException(err)
		return nil, err
	}

	return users_response, nil
}

func (usersGrpc *MockUsersGrpc) GetUserById(id int) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Id:                 1,
		Username:           "testuser",
		Phone:              "testphone",
		Email:              "testemail",
		FirstName:          "testfirstname",
		LastName:           "testlastname",
		MiddleName:         "testmiddlename",
		Activity:           "testactivity",
		Status:             "teststatus",
		EmailConfirmed:     true,
		PhoneConfirmed:     true,
		LastSeen:           "2023-01-01T20:00:00Z",
		OriginalAvatarUrl:  "testurl",
		ConvertedAvatarUrl: "testurl",
	}, nil
}

func (usersGrpc *MockUsersGrpc) GetUsersByIds(ids []int) (*pb.UsersArrayResponse, error) {
	var users_array []*pb.UserResponse
	users_array = append(users_array, &pb.UserResponse{
		Id:                 1,
		Username:           "testuser",
		Phone:              "testphone",
		Email:              "testemail",
		FirstName:          "testfirstname",
		LastName:           "testlastname",
		MiddleName:         "testmiddlename",
		Activity:           "testactivity",
		Status:             "teststatus",
		EmailConfirmed:     true,
		PhoneConfirmed:     true,
		LastSeen:           "2023-01-01T20:00:00Z",
		OriginalAvatarUrl:  "testurl",
		ConvertedAvatarUrl: "testurl",
	})
	return &pb.UsersArrayResponse{
		Users: users_array,
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
