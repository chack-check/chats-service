package redisdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/entities"

	"github.com/redis/go-redis/v9"
)

type RedisActionUser struct {
	ID         int
	LastName   string
	FirstName  string
	MiddleName *string
	Username   string
}

type UserActionsLoggingAdapter struct {
	repository ports.UserActionsRepositoryPort
}

func (adapter UserActionsLoggingAdapter) AddChatActionUser(chat entities.Chat, user entities.User, actionType constants.ActionTypes) map[constants.ActionTypes][]entities.ActionUser {
	log.Printf("adding chat action user: chat=%+v, user=%+v, actionType=%v", chat, user, actionType)
	actions := adapter.repository.AddChatActionUser(chat, user, actionType)
	log.Printf("chat actions: %+v", actions)
	return actions
}

func (adapter UserActionsLoggingAdapter) RemoveChatActionUser(chat entities.Chat, userID int, actionType constants.ActionTypes) map[constants.ActionTypes][]entities.ActionUser {
	log.Printf("removing chat action user: chat=%+v, userId=%d, actionType=%v", chat, userID, actionType)
	actions := adapter.repository.RemoveChatActionUser(chat, userID, actionType)
	log.Printf("chat actions: %+v", actions)
	return actions
}

func (adapter UserActionsLoggingAdapter) GetAllChatActionsUsers(chat entities.Chat) map[constants.ActionTypes][]entities.ActionUser {
	log.Printf("fetching all chat actions users: chat=%+v", chat)
	actions := adapter.repository.GetAllChatActionsUsers(chat)
	log.Printf("chat actions: %+v", actions)
	return actions
}

type UserActionsAdapter struct {
	db *redis.Client
}

func (adapter UserActionsAdapter) getChatActionsKey(chatID int) string {
	return fmt.Sprintf("chat:%d:actions", chatID)
}

func (adapter UserActionsAdapter) AddChatActionUser(chat entities.Chat, user entities.User, actionType constants.ActionTypes) map[constants.ActionTypes][]entities.ActionUser {
	actionUsers, err := adapter.db.HGet(context.Background(), adapter.getChatActionsKey(chat.GetID()), string(actionType)).Result()
	log.Printf("action users: %+v, err: %v", actionUsers, err)
	if err != nil && err != redis.Nil {
		return adapter.GetAllChatActionsUsers(chat)
	}
	if actionUsers == "" || actionUsers == "null" {
		actionUsers = "[]"
	}

	users := make([]RedisActionUser, 0, 1)
	err = json.Unmarshal([]byte(actionUsers), &users)
	log.Printf("users: %+v, err: %v", users, err)
	if err != nil {
		return adapter.GetAllChatActionsUsers(chat)
	}

	users = append(users, RedisActionUser{
		ID:         user.GetID(),
		LastName:   user.GetLastName(),
		FirstName:  user.GetFirstName(),
		MiddleName: user.GetMiddleName(),
		Username:   user.GetUsername(),
	})
	usersJSON, err := json.Marshal(users)
	log.Printf("users: %+v, usersJson: %s, err: %v", users, usersJSON, err)
	if err != nil {
		return adapter.GetAllChatActionsUsers(chat)
	}

	adapter.db.HSet(context.Background(), adapter.getChatActionsKey(chat.GetID()), string(actionType), string(usersJSON)).Result()
	return adapter.GetAllChatActionsUsers(chat)
}

func (adapter UserActionsAdapter) RemoveChatActionUser(chat entities.Chat, userID int, actionType constants.ActionTypes) map[constants.ActionTypes][]entities.ActionUser {
	actionUsers, err := adapter.db.HGet(context.Background(), adapter.getChatActionsKey(chat.GetID()), string(actionType)).Result()
	if err != nil || actionUsers == "" {
		return adapter.GetAllChatActionsUsers(chat)
	}

	var users []RedisActionUser
	err = json.Unmarshal([]byte(actionUsers), &users)
	if err != nil {
		return adapter.GetAllChatActionsUsers(chat)
	}

	var resultUsers []RedisActionUser
	for _, user := range users {
		if user.ID == userID {
			continue
		}

		resultUsers = append(resultUsers, user)
	}

	usersJSON, err := json.Marshal(resultUsers)
	if err != nil {
		return adapter.GetAllChatActionsUsers(chat)
	}

	adapter.db.HSet(context.Background(), adapter.getChatActionsKey(chat.GetID()), string(actionType), string(usersJSON))
	return adapter.GetAllChatActionsUsers(chat)
}

func (adapter UserActionsAdapter) GetAllChatActionsUsers(chat entities.Chat) map[constants.ActionTypes][]entities.ActionUser {
	chatAllActions, err := adapter.db.HGetAll(context.Background(), adapter.getChatActionsKey(chat.GetID())).Result()
	if err != nil {
		return map[constants.ActionTypes][]entities.ActionUser{}
	}

	actions := make(map[constants.ActionTypes][]entities.ActionUser)
	for key, value := range chatAllActions {
		actionType := constants.ActionTypes(key)
		var actionUsers []RedisActionUser
		err = json.Unmarshal([]byte(value), &actionUsers)
		if err != nil {
			continue
		}

		var actionUsersSchemas []entities.ActionUser
		for _, user := range actionUsers {
			actionUsersSchemas = append(actionUsersSchemas, entities.NewActionUser(user.ID, user.LastName, user.FirstName, user.MiddleName, user.Username))
		}

		actions[actionType] = actionUsersSchemas
	}

	return actions
}

func NewUserActionsAdapter(db *redis.Client) ports.UserActionsRepositoryPort {
	return UserActionsLoggingAdapter{repository: UserActionsAdapter{db: db}}
}
