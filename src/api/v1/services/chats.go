package services

import (
	"fmt"
	"log"
	"slices"

	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/api/v1/schemas"
	"github.com/chack-check/chats-service/grpc_client"
	"github.com/chack-check/chats-service/protousers"
	"github.com/chack-check/chats-service/rabbit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

func setupChatTitleAndAvatar(chat *models.Chat, chatUser *protousers.UserResponse) {
	chat.Title = fmt.Sprintf("%s %s", chatUser.LastName, chatUser.FirstName)
	if chatUser.ConvertedAvatarUrl != "" {
		chat.AvatarURL = chatUser.ConvertedAvatarUrl
	} else {
		chat.AvatarURL = chatUser.OriginalAvatarUrl
	}
}

func setupUserChatData(chat *models.Chat, currentUserId int) {
	var anotherUserId int64
	for _, member := range chat.Members {
		if member != int64(currentUserId) {
			anotherUserId = member
		}
	}

	anotherUser, err := grpc_client.UsersGrpcClient.GetUserById(int(anotherUserId))
	if err != nil {
		chat.Title = "Untitled"
		chat.AvatarURL = ""
	} else {
		setupChatTitleAndAvatar(chat, anotherUser)
	}
}

func setupUserManyChatsData(chats []*models.Chat, currentUserId int) {
	var anotherUsersIds []int
	for _, chat := range chats {
		if chat.Type != "user" {
			continue
		}

		for _, member := range chat.Members {
			if member != int64(currentUserId) && !slices.Contains(anotherUsersIds, int(member)) {
				anotherUsersIds = append(anotherUsersIds, int(member))
			}
		}
	}

	anotherUsers, err := grpc_client.UsersGrpcClient.GetUsersByIds(anotherUsersIds)
	if err != nil {
		return
	}

	for _, chat := range chats {
		if chat.Type != "user" {
			continue
		}

		var anotherUser *protousers.UserResponse

		for _, member := range chat.Members {
			if member == int64(currentUserId) {
				continue
			}

			for _, user := range anotherUsers.Users {
				if user.Id == int32(member) {
					anotherUser = user
				}
			}
		}

		setupChatTitleAndAvatar(chat, anotherUser)
	}
}

type ChatsManager struct {
	ChatsQueries *models.ChatsQueries
}

func (manager *ChatsManager) GetConcrete(chatID uint, token *jwt.Token) (*models.Chat, error) {
	tokenSubject, err := GetTokenSubject(token)

	if err != nil {
		return nil, err
	}

	chat, err := manager.ChatsQueries.GetWithMember(chatID, uint(tokenSubject.UserId))
	if err != nil {
		return nil, err
	}

	if chat.Type == "user" {
		setupUserChatData(chat, tokenSubject.UserId)
	}

	return chat, nil
}

func (manager *ChatsManager) GetByIds(chatsIds []int, token *jwt.Token) ([]*models.Chat, error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	chats, err := manager.ChatsQueries.GetByIds(chatsIds, tokenSubject.UserId)
	if err != nil {
		return nil, err
	}

	setupUserManyChatsData(chats, tokenSubject.UserId)

	return chats, nil
}

func (manager *ChatsManager) GetAll(token *jwt.Token, page int, perPage int) *schemas.PaginatedResponse[*models.Chat] {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil
	}

	count := manager.ChatsQueries.GetAllWithMemberCount(uint(tokenSubject.UserId))
	countValue := *count
	chats := manager.ChatsQueries.GetAllWithMember(uint(tokenSubject.UserId), page, perPage)

	setupUserManyChatsData(chats, tokenSubject.UserId)

	paginatedResponse := schemas.NewPaginatedResponse(page, perPage, int(countValue), chats)
	return &paginatedResponse
}

func (manager *ChatsManager) createGroupChat(chat *models.Chat, userId int) error {
	chat.OwnerId = uint(userId)
	chat.Type = "group"

	if err := manager.ChatsQueries.Create(chat); err != nil {
		return err
	}

	return nil
}

func (manager *ChatsManager) validateUserChatExists(userId int, anotherUserId int) error {
	alreadyExists := manager.ChatsQueries.GetExistingWithUser(uint(userId), uint(anotherUserId))
	if alreadyExists {
		return fmt.Errorf("you already have chat with this user")
	}

	return nil
}

func (manager *ChatsManager) createUserChat(chat *models.Chat, currentUser *protousers.UserResponse, chatUser *protousers.UserResponse) error {
	chat.Members = []int64{int64(currentUser.Id), int64(chatUser.Id)}
	chat.OwnerId = 0
	chat.Title = ""
	chat.Type = "user"
	chat.AvatarURL = ""

	err := manager.validateUserChatExists(int(currentUser.Id), int(chatUser.Id))
	if err != nil {
		return err
	}

	log.Printf("Creating chat: %v", chat)

	if err := manager.ChatsQueries.Create(chat); err != nil {
		return err
	}

	setupChatTitleAndAvatar(chat, chatUser)
	return nil
}

func (manager *ChatsManager) sendChatEvent(chat *models.Chat, eventType string) error {
	var included_users []int
	for _, user := range chat.Members {
		included_users = append(included_users, int(user))
	}

	chatEvent, err := rabbit.NewSystemEvent(eventType, included_users, chat)
	if err != nil {
		return err
	}

	err = rabbit.EventsRabbitConnection.SendEvent(chatEvent)
	log.Printf("Sended event with type %s to rabbitmq", eventType)

	if err != nil {
		log.Printf("Error when publishing event with type %s in queue: %v", eventType, err)
		return fmt.Errorf("error sending event with type %s", eventType)
	}

	return nil
}

func (manager *ChatsManager) Search(query string, token *jwt.Token, page int, perPage int) (*schemas.PaginatedResponse[*models.Chat], error) {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil, err
	}

	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 1
	}

	chatsCount := manager.ChatsQueries.SearchCount(uint(tokenSubject.UserId), query, page, perPage)
	chats := manager.ChatsQueries.Search(uint(tokenSubject.UserId), query, page, perPage)

	setupUserManyChatsData(chats, tokenSubject.UserId)

	response := schemas.NewPaginatedResponse(page, perPage, int(chatsCount), chats)
	return &response, nil
}

func (manager *ChatsManager) Create(chat *models.Chat, token *jwt.Token, chatUserId uint) error {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return err
	}

	if chatUserId == 0 {
		err = manager.createGroupChat(chat, tokenSubject.UserId)
		if err != nil {
			return err
		}

		manager.sendChatEvent(chat, "chat_created")
		return nil
	}

	if chatUserId == uint(tokenSubject.UserId) {
		return fmt.Errorf("you can't create chat for you")
	}

	members := pq.Int32Array{int32(tokenSubject.UserId), int32(chatUserId)}
	if chatId := manager.ChatsQueries.GetDeletedChatId(members); chatId > 0 {
		manager.ChatsQueries.RestoreChat(chatId)
		if err := manager.sendChatEvent(chat, "chat_created"); err != nil {
			return err
		}

		return nil
	}

	currentUser, err := grpc_client.UsersGrpcClient.GetUserById(tokenSubject.UserId)

	if err != nil || currentUser == nil {
		return fmt.Errorf("there is no user with id %d", tokenSubject.UserId)
	}

	chatUser, err := grpc_client.UsersGrpcClient.GetUserById(int(chatUserId))

	if err != nil || chatUser == nil {
		return fmt.Errorf("there is no user with id %d", chatUserId)
	}

	log.Printf("Creating user chat: %v, user id: %v, chat user: %v", chat, tokenSubject.UserId, chatUser)
	err = manager.createUserChat(chat, currentUser, chatUser)
	if err != nil {
		return err
	}

	manager.sendChatEvent(chat, "chat_created")
	return nil
}

func (manager *ChatsManager) Delete(chat *models.Chat) {
	manager.ChatsQueries.Delete(chat)
	manager.sendChatEvent(chat, "chat_deleted")
}

func NewChatsManager() *ChatsManager {
	return &ChatsManager{
		ChatsQueries: &models.ChatsQueries{},
	}
}
