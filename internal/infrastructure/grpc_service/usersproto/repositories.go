package usersproto

import (
	"context"
	"fmt"
	"log"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
	"chats-service/internal/infrastructure/grpc_service/usersproto/usersprotobuf"
)

var (
	ErrUserNotFound = fmt.Errorf("user not found")
)

type UsersLoggingRepository struct {
	repository ports.UsersRepository
}

func (repository UsersLoggingRepository) GetById(id int) (*entities.User, error) {
	log.Printf("fetching user by id: %d", id)
	user, err := repository.repository.GetById(id)
	if err != nil {
		log.Printf("error fetching user by id: %v", err)
		return user, err
	}

	log.Printf("fetched user: %+v", user)
	return user, err
}

func (repository UsersLoggingRepository) GetByIds(ids []int) []entities.User {
	log.Printf("fetching users by ids: %v", ids)
	users := repository.repository.GetByIds(ids)
	log.Printf("fetched users: %+v", users)
	return users
}

type UsersRepository struct {
	client usersprotobuf.UsersClient
}

func (repository UsersRepository) GetById(id int) (*entities.User, error) {
	user, err := repository.client.GetUserById(context.Background(), &usersprotobuf.GetUserByIdRequest{Id: int32(id)})
	if err != nil {
		log.Printf("error finding user by id %d: %v", id, err)
		return nil, ErrUserNotFound
	}

	userModel := ProtoUserToModel(user)
	return &userModel, nil
}

func (repository UsersRepository) GetByIds(ids []int) []entities.User {
	userIds := make([]int32, 0, len(ids))
	for _, id := range ids {
		userIds = append(userIds, int32(id))
	}

	foundedUsers, err := repository.client.GetUsersByIds(context.Background(), &usersprotobuf.GetUsersByIdsRequest{Ids: userIds})
	var usersModels []entities.User
	if err != nil {
		return usersModels
	}

	for _, user := range foundedUsers.Users {
		userModel := ProtoUserToModel(user)
		usersModels = append(usersModels, userModel)
	}

	return usersModels
}

func NewUsersRepository(client usersprotobuf.UsersClient) ports.UsersRepository {
	return UsersLoggingRepository{repository: UsersRepository{client: client}}
}
