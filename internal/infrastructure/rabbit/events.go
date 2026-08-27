package rabbit

import (
	"encoding/json"
	"log"

	"chats-service/internal/application/controllers"
	"chats-service/internal/domain/constants"
	"chats-service/internal/domain/dtos"
	"chats-service/internal/infrastructure/database"
)

type SystemEvent struct {
	IncludedUsers []int  `json:"included_users"`
	EventType     string `json:"event_type"`
	Data          string `json:"data"`
}

func NewSystemEvent(eventType string, includedUsers []int, data interface{}) (*SystemEvent, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &SystemEvent{IncludedUsers: includedUsers, EventType: eventType, Data: string(jsonData)}, nil
}

type RecognitionEvent struct {
	MessageID int    `json:"message_id"`
	Content   string `json:"content"`
}

func HandleUserCreated(event SystemEvent) {
	var eventUser EventUser
	err := json.Unmarshal([]byte(event.Data), &eventUser)
	if err != nil {
		log.Printf("error unmarshaling event user data: %v", err)
	}

	data := dtos.NewCreateChatData(constants.SavedMessagesChatType, nil, nil, []int{}, &eventUser.Id)
	controller := controllers.NewCreateSavedMessagesChatController(
		database.NewChatsAdapter(*database.DatabaseConnection),
	)
	controller.Execute(data, eventUser.Id)
}

func HandleMessageRecognized(messageId int, content string) {
	controller := controllers.NewRecognizeMessageController(
		database.NewMessagesAdapter(*database.DatabaseConnection),
		NewMessageEventsAdapter(*EventsRabbitConnection),
	)
	controller.Execute(messageId, content)
}
