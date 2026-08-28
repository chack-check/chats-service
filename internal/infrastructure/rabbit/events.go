package rabbit

import (
	"encoding/json"
	"log"

	"chats-service/internal/application/dtos"
	"chats-service/internal/application/ports"
	"chats-service/internal/application/usecases"
	"chats-service/internal/domain/constants"
	"chats-service/internal/infrastructure/database"
	"chats-service/internal/infrastructure/database/repositories"
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
	useCase := usecases.NewCreateSavedMessagesChatUseCase(
		repositories.NewChatsRepository(*database.DatabaseConnection),
	)
	useCase.Execute(data, eventUser.Id)
}

func HandleMessageRecognized(messageId int, content string, publisher ports.MessageEventsPublisher) {
	useCase := usecases.NewRecognizeMessageUseCase(
		repositories.NewMessagesRepository(*database.DatabaseConnection),
		publisher,
	)
	useCase.Execute(messageId, content)
}
