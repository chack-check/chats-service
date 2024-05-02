package usersproto

import (
	"context"
	"fmt"
	"log"

	"github.com/chack-check/chats-service/domain/users"
	"github.com/chack-check/chats-service/infrastructure/grpc_service/usersproto/usersprotobuf"
)

var (
	ErrUserNotFound = fmt.Errorf("user not found")
)

type UsersLoggingAdapter struct {
	adapter users.UsersPort
}

func (adapter UsersLoggingAdapter) GetById(id int) (*users.User, error) {
	log.Printf("fetching user by id: %d", id)
	user, err := adapter.adapter.GetById(id)
	if err != nil {
		log.Printf("error fetching user by id: %v", err)
		return user, err
	}

	log.Printf("fetched user: %+v", user)
	return user, err
}

func (adapter UsersLoggingAdapter) GetByIds(ids []int) []users.User {
	log.Printf("fetching users by ids: %v", ids)
	users := adapter.adapter.GetByIds(ids)
	log.Printf("fetched users: %+v", users)
	return users
}

type UsersAdapter struct {
	client usersprotobuf.UsersClient
}

func (adapter UsersAdapter) GetById(id int) (*users.User, error) {
	user, err := adapter.client.GetUserById(context.Background(), &usersprotobuf.GetUserByIdRequest{Id: int32(id)})
	if err != nil {
		log.Printf("error finding user by id %d: %v", id, err)
		return nil, ErrUserNotFound
	}

	userModel := ProtoUserToModel(user)
	return &userModel, nil
}

func (adapter UsersAdapter) GetByIds(ids []int) []users.User {
	var userIds []int32
	for _, id := range ids {
		userIds = append(userIds, int32(id))
	}

	foundedUsers, err := adapter.client.GetUsersByIds(context.Background(), &usersprotobuf.GetUsersByIdsRequest{Ids: userIds})
	var usersModels []users.User
	if err != nil {
		return usersModels
	}

	for _, user := range foundedUsers.Users {
		userModel := ProtoUserToModel(user)
		usersModels = append(usersModels, userModel)
	}

	return usersModels
}

func NewUsersAdapter(client usersprotobuf.UsersClient) users.UsersPort {
	return UsersLoggingAdapter{adapter: UsersAdapter{client: client}}
}
