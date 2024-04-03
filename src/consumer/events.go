package consumer

import (
	"encoding/json"
	"log"

	"github.com/chack-check/chats-service/api/v1/dtos"
	"github.com/chack-check/chats-service/api/v1/services"
	"github.com/getsentry/sentry-go"
)

type DataStruct struct {
	EventType     string `json:"event_type"`
	IncludedUsers []int  `json:"included_users"`
	Data          string `json:"data"`
}

type EventUser struct {
	Id         int     `json:"id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	MiddleName *string `json:"middle_name"`
	Username   string  `json:"username"`
}

func HandleUserCreated(message []byte) {
	log.Printf("Handling create user event: %+v", message)
	var data DataStruct
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Printf("Error when unmarshalling user created event: %v", err)
		sentry.CaptureException(err)
		return
	}
	var userEvent EventUser

	err = json.Unmarshal([]byte(data.Data), &userEvent)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Error when unmarshalling user created event: %v", err)
		return
	}

	chatsService := services.NewChatsManager()
	members := []int{userEvent.Id}
	chatDto := dtos.ChatDto{
		Title:   "Saved messages",
		Type:    "saved_messages",
		Members: members,
		OwnerId: userEvent.Id,
	}
	chatsService.SystemSave(&chatDto)
}
