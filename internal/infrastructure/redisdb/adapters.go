package redisdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"chats-service/internal/domain/chats"
	"chats-service/internal/domain/users"
	"github.com/redis/go-redis/v9"
)

type RedisActionUser struct {
	Id         int
	LastName   string
	FirstName  string
	MiddleName *string
	Username   string
}

type UserActionsLoggingAdapter struct {
	adapter chats.UserActionsPort
}

func (adapter UserActionsLoggingAdapter) AddChatActionUser(chat chats.Chat, user users.User, actionType chats.ActionTypes) map[chats.ActionTypes][]users.ActionUser {
	log.Printf("adding chat action user: chat=%+v, user=%+v, actionType=%v", chat, user, actionType)
	actions := adapter.adapter.AddChatActionUser(chat, user, actionType)
	log.Printf("chat actions: %+v", actions)
	return actions
}

func (adapter UserActionsLoggingAdapter) RemoveChatActionUser(chat chats.Chat, userID int, actionType chats.ActionTypes) map[chats.ActionTypes][]users.ActionUser {
	log.Printf("removing chat action user: chat=%+v, userId=%d, actionType=%v", chat, userID, actionType)
	actions := adapter.adapter.RemoveChatActionUser(chat, userID, actionType)
	log.Printf("chat actions: %+v", actions)
	return actions
}

func (adapter UserActionsLoggingAdapter) GetAllChatActionsUsers(chat chats.Chat) map[chats.ActionTypes][]users.ActionUser {
	log.Printf("fetching all chat actions users: chat=%+v", chat)
	actions := adapter.adapter.GetAllChatActionsUsers(chat)
	log.Printf("chat actions: %+v", actions)
	return actions
}

type UserActionsAdapter struct {
	db *redis.Client
}

func (adapter UserActionsAdapter) getChatActionsKey(chatID int) string {
	return fmt.Sprintf("chat:%d:actions", chatID)
}

func (adapter UserActionsAdapter) AddChatActionUser(chat chats.Chat, user users.User, actionType chats.ActionTypes) map[chats.ActionTypes][]users.ActionUser {
	actionUsers, err := adapter.db.HGet(context.Background(), adapter.getChatActionsKey(chat.GetId()), string(actionType)).Result()
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
		Id:         user.GetId(),
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

	adapter.db.HSet(context.Background(), adapter.getChatActionsKey(chat.GetId()), string(actionType), string(usersJSON)).Result()
	return adapter.GetAllChatActionsUsers(chat)
}

func (adapter UserActionsAdapter) RemoveChatActionUser(chat chats.Chat, userID int, actionType chats.ActionTypes) map[chats.ActionTypes][]users.ActionUser {
	actionUsers, err := adapter.db.HGet(context.Background(), adapter.getChatActionsKey(chat.GetId()), string(actionType)).Result()
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
		if user.Id == userID {
			continue
		}

		resultUsers = append(resultUsers, user)
	}

	usersJSON, err := json.Marshal(resultUsers)
	if err != nil {
		return adapter.GetAllChatActionsUsers(chat)
	}

	adapter.db.HSet(context.Background(), adapter.getChatActionsKey(chat.GetId()), string(actionType), string(usersJSON))
	return adapter.GetAllChatActionsUsers(chat)
}

func (adapter UserActionsAdapter) GetAllChatActionsUsers(chat chats.Chat) map[chats.ActionTypes][]users.ActionUser {
	chatAllActions, err := adapter.db.HGetAll(context.Background(), adapter.getChatActionsKey(chat.GetId())).Result()
	if err != nil {
		return map[chats.ActionTypes][]users.ActionUser{}
	}

	actions := make(map[chats.ActionTypes][]users.ActionUser)
	for key, value := range chatAllActions {
		actionType := chats.ActionTypes(key)
		var actionUsers []RedisActionUser
		err = json.Unmarshal([]byte(value), &actionUsers)
		if err != nil {
			continue
		}

		var actionUsersSchemas []users.ActionUser
		for _, user := range actionUsers {
			actionUsersSchemas = append(actionUsersSchemas, users.NewActionUser(user.Id, user.LastName, user.FirstName, user.MiddleName, user.Username))
		}

		actions[actionType] = actionUsersSchemas
	}

	return actions
}

func NewUserActionsAdapter(db *redis.Client) chats.UserActionsPort {
	return UserActionsLoggingAdapter{adapter: UserActionsAdapter{db: db}}
}
