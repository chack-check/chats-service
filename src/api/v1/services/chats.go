package services

import (
	"fmt"

	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/api/v1/schemas"
	"github.com/chack-check/chats-service/grpc_client"
	"github.com/chack-check/chats-service/protousers"
)

type ChatsManager struct {
	ChatsQueries *models.ChatsQueries
}

func (manager *ChatsManager) GetConcrete(chatID uint, user *protousers.UserResponse) (*models.Chat, error) {
	chat, err := manager.ChatsQueries.GetWithMember(chatID, uint(user.Id))

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (manager *ChatsManager) GetAll(user *protousers.UserResponse, page *int, perPage *int) *schemas.PaginatedResponse[models.Chat] {
	count := manager.ChatsQueries.GetAllWithMemberCount(uint(user.Id))
	countValue := *count
	chats := manager.ChatsQueries.GetAllWithMember(uint(user.Id), page, perPage)
	paginatedResponse := schemas.NewPaginatedResponse[models.Chat](*page, *perPage, int(countValue), *chats)
	return &paginatedResponse
}

func (manager *ChatsManager) createGroupChat(chat *models.Chat, user *protousers.UserResponse) error {
	chat.OwnerId = uint(user.Id)
	chat.Type = "group"

	if err := manager.ChatsQueries.Create(chat); err != nil {
		return err
	}

	return nil
}

func (manager *ChatsManager) createUserChat(chat *models.Chat, user *protousers.UserResponse, chatUser *protousers.UserResponse) error {
	chat.Members = []int64{int64(user.Id), int64(chatUser.Id)}
	chat.OwnerId = 0
	chat.Title = fmt.Sprintf("%v %v %v", chatUser.LastName, chatUser.FirstName, chatUser.MiddleName)
	chat.Type = "user"
	chat.AvatarURL = "https://google.com"

	if err := manager.ChatsQueries.Create(chat); err != nil {
		return err
	}

	return nil
}

func (manager *ChatsManager) Create(chat *models.Chat, user *protousers.UserResponse, chatUserId uint) error {
	if chatUserId == 0 {
		return manager.createGroupChat(chat, user)
	}

	if chatUserId == uint(user.Id) {
		return fmt.Errorf("You can't create chat for you")
	}

	chatUser, err := grpc_client.UsersGrpcClient.GetUserById(int(chatUserId))

	if err != nil {
		return fmt.Errorf("There is no user with id %d", chatUserId)
	}

	return manager.createUserChat(chat, user, chatUser)
}

func NewChatsManager() *ChatsManager {
	return &ChatsManager{
		ChatsQueries: &models.ChatsQueries{},
	}
}
