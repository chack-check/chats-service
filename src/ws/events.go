package ws

import (
	"encoding/json"

	"github.com/chack-check/chats-service/api/v1/models"
	"github.com/chack-check/chats-service/grpc_client"
)

type ClientMessageHandler interface {
	HandleMessage(message []byte, client *Client) ([]byte, error)
}

type AuthenticationEvent struct {
	Type         string `json:"type"`
	RefreshToken string `json:"refreshToken"`
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type AuthMessageHandler struct{}

func (handler AuthMessageHandler) authenticateUser(authMessage *AuthenticationEvent, client *Client) error {
	user, err := grpc_client.UsersGrpcClient.GetUserByRefreshToken(authMessage.RefreshToken)

	if err != nil || user == nil {
		return err
	}

	client.user = user
	return nil
}

func (handler AuthMessageHandler) HandleMessage(message []byte, client *Client) ([]byte, error) {
	authEvent := AuthenticationEvent{}
	authEventPointer := &authEvent
	err := json.Unmarshal(message, authEventPointer)

	if err != nil {
		return []byte{}, nil
	}

	err = handler.authenticateUser(authEventPointer, client)

	if err != nil {
		return NewAuthenticationErrorEventBytes(), err
	}

	return NewAuthenticationSuccessEventBytes(), nil
}

func NewAuthenticationErrorEventBytes() []byte {
	responseEvent := Event{Type: "AuthenticationError"}
	responseEventJson, _ := json.Marshal(responseEvent)
	return responseEventJson
}

func NewAuthenticationSuccessEventBytes() []byte {
	responseEvent := Event{Type: "AuthenticationSuccess"}
	responseEventJson, _ := json.Marshal(responseEvent)
	return responseEventJson
}

func NewUndefinedMessageBytes() []byte {
	responseMessage := Event{Type: "UndefinedMessage"}
	responseMessageJson, _ := json.Marshal(responseMessage)
	return responseMessageJson
}

func NewMessageEvent(message *models.Message) []byte {
	messageEvent := Event{Type: "NewMessage", Data: *message}
	messageEventJson, _ := json.Marshal(messageEvent)
	return messageEventJson
}
