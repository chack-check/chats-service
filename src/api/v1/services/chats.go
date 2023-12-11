package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/api/v1/schemas"
	"github.com/chack-check/chats-service/grpc_client"
	"github.com/chack-check/chats-service/protousers"
	"github.com/chack-check/chats-service/rabbit"
	"github.com/golang-jwt/jwt/v5"
)

func getChatEventFromChat(chat *models.Chat) *rabbit.ChatEvent {
	members := []int{}
	admins := []int{}

	for _, member := range chat.Members {
		members = append(members, int(member))
	}

	for _, admin := range chat.Admins {
		admins = append(admins, int(admin))
	}

	return &rabbit.ChatEvent{
		Type:          "chat",
		IncludedUsers: members,
		ChatId:        int(chat.ID),
		AvatarURL:     chat.AvatarURL,
		Title:         chat.Title,
		ChatType:      chat.Type,
		Members:       members,
		IsArchived:    chat.IsArchived,
		OwnerID:       int(chat.OwnerId),
		Admins:        admins,
	}
}

type ChatsManager struct {
	ChatsQueries *models.ChatsQueries
}

func (manager *ChatsManager) decodeUserChatTitle(currentUserId int, title string) string {
	decodedTitles := map[string]int{}
	err := json.Unmarshal([]byte(title), &decodedTitles)

	if err != nil {
		return "Untitled"
	}

	for chatTitle, id := range decodedTitles {
		if id != currentUserId {
			return chatTitle
		}
	}

	return "Untitled"
}

func (manager *ChatsManager) decodeUserChatAvatar(currentUserId int, avatarUrl string) string {
	decodedAvatars := map[string]int{}
	err := json.Unmarshal([]byte(avatarUrl), &decodedAvatars)

	if err != nil {
		return "Untitled"
	}

	for chatAvatarUrl, id := range decodedAvatars {
		if id != currentUserId {
			return chatAvatarUrl
		}
	}

	return "Untitled"
}

func (manager *ChatsManager) GetConcrete(chatID uint, token *jwt.Token) (*models.Chat, error) {
	tokenSubject, err := GetTokenSubject(token)

	if err != nil {
		return nil, err
	}

	chat, err := manager.ChatsQueries.GetWithMember(chatID, uint(tokenSubject.UserId))
	chat.Title = manager.decodeUserChatTitle(tokenSubject.UserId, chat.Title)
	chat.AvatarURL = manager.decodeUserChatAvatar(tokenSubject.UserId, chat.AvatarURL)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (manager *ChatsManager) GetAll(token *jwt.Token, page int, perPage int) *schemas.PaginatedResponse[models.Chat] {
	tokenSubject, err := GetTokenSubject(token)
	if err != nil {
		return nil
	}

	count := manager.ChatsQueries.GetAllWithMemberCount(uint(tokenSubject.UserId))
	countValue := *count
	chats := manager.ChatsQueries.GetAllWithMember(uint(tokenSubject.UserId), page, perPage)
	paginatedResponse := schemas.NewPaginatedResponse(page, perPage, int(countValue), *chats)
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
		return fmt.Errorf("You already have chat with this user")
	}

	return nil
}

func (manager *ChatsManager) createUserChat(chat *models.Chat, currentUser *protousers.UserResponse, chatUser *protousers.UserResponse) error {
	chatTitleForCurrentUser := fmt.Sprintf("%s %s", currentUser.LastName, currentUser.FirstName)
	chatTitleForChatUser := fmt.Sprintf("%s %s", chatUser.LastName, chatUser.FirstName)
	chatTitles := map[string]int{
		chatTitleForCurrentUser: int(currentUser.Id),
		chatTitleForChatUser:    int(chatUser.Id),
	}
	chatAvatarUrls := map[string]int{
		currentUser.AvatarUrl: int(currentUser.Id),
		chatUser.AvatarUrl:    int(chatUser.Id),
	}
	jsonTitles, err := json.Marshal(chatTitles)
	if err != nil {
		return err
	}

	jsonAvatars, err := json.Marshal(chatAvatarUrls)
	if err != nil {
		return err
	}

	chat.Members = []int64{int64(currentUser.Id), int64(chatUser.Id)}
	chat.OwnerId = 0
	chat.Title = string(jsonTitles)
	chat.Type = "user"
	chat.AvatarURL = string(jsonAvatars)

	err = manager.validateUserChatExists(int(currentUser.Id), int(chatUser.Id))
	if err != nil {
		return err
	}

	log.Printf("Creating chat: %v", chat)

	if err := manager.ChatsQueries.Create(chat); err != nil {
		return err
	}

	return nil
}

func (manager *ChatsManager) sendChatEvent(chat *models.Chat) error {
	chatEvent := getChatEventFromChat(chat)
	err := rabbit.EventsRabbitConnection.SendEvent(chatEvent)
	log.Printf("Sended chat to rabbitmq")

	if err != nil {
		log.Printf("Error when publishing chat event in queue: %v", err)
		return fmt.Errorf("Error sending chat event")
	}

	return nil
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

		manager.sendChatEvent(chat)
		return nil
	}

	if chatUserId == uint(tokenSubject.UserId) {
		return fmt.Errorf("You can't create chat for you")
	}

	currentUser, err := grpc_client.UsersGrpcClient.GetUserById(tokenSubject.UserId)

	if err != nil || currentUser == nil {
		return fmt.Errorf("There is no user with id %d", tokenSubject.UserId)
	}

	chatUser, err := grpc_client.UsersGrpcClient.GetUserById(int(chatUserId))

	if err != nil || chatUser == nil {
		return fmt.Errorf("There is no user with id %d", chatUserId)
	}

	log.Printf("Creating user chat: %v, user id: %v, chat user: %v", chat, tokenSubject.UserId, chatUser)
	err = manager.createUserChat(chat, currentUser, chatUser)
	if err != nil {
		return err
	}

	manager.sendChatEvent(chat)
	return nil
}

func NewChatsManager() *ChatsManager {
	return &ChatsManager{
		ChatsQueries: &models.ChatsQueries{},
	}
}
