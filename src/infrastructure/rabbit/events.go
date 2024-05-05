package rabbit

import (
	"encoding/json"
	"log"

	"github.com/chack-check/chats-service/domain/chats"
	"github.com/chack-check/chats-service/infrastructure/database"
)

type SystemEvent struct {
	IncludedUsers []int  `json:"included_users"`
	EventType     string `json:"event_type"`
	Data          string `json:"data"`
}

func NewSystemEvent(eventType string, includedUsers []int, data interface{}) (*SystemEvent, error) {
	json_data, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &SystemEvent{IncludedUsers: includedUsers, EventType: eventType, Data: string(json_data)}, nil
}

func HandleUserCreated(event SystemEvent) {
	var eventUser EventUser
	err := json.Unmarshal([]byte(event.Data), &eventUser)
	if err != nil {
		log.Printf("error unmarshaling event user data: %v", err)
	}

	data := chats.NewCreateChatData(chats.SavedMessagesChatType, nil, nil, []int{}, &eventUser.Id)
	handler := chats.NewCreateSavedMessagesChatHandler(
		database.NewChatsAdapter(*database.DatabaseConnection),
	)
	handler.Execute(data, eventUser.Id)
}
