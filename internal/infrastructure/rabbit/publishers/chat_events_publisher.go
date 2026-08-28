package publishers

import (
	"log"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
	"chats-service/internal/infrastructure/rabbit"
)

type ChatEventsPublisher struct {
	connection rabbit.RabbitConnection
}

func (publisher ChatEventsPublisher) getSystemEventForChat(chat entities.Chat, eventType string) (*rabbit.SystemEvent, error) {
	chatEvent := rabbit.ChatToChatEvent(chat)
	return rabbit.NewSystemEvent(eventType, chat.GetMembers(), chatEvent)
}

func (publisher ChatEventsPublisher) sendChatEvent(chat entities.Chat, eventType string) {
	systemEvent, err := publisher.getSystemEventForChat(chat, eventType)
	if err != nil {
		return
	}

	publisher.connection.SendEvent(systemEvent)
}

func (publisher ChatEventsPublisher) SendChatCreated(chat entities.Chat) {
	publisher.sendChatEvent(chat, "chat_created")
}

func (publisher ChatEventsPublisher) SendChatDeleted(chat entities.Chat) {
	publisher.sendChatEvent(chat, "chat_deleted")
}

func (publisher ChatEventsPublisher) SendChatUserAction(chat entities.Chat) {
	publisher.sendChatEvent(chat, "chat_user_action")
}

func (publisher ChatEventsPublisher) SendChatChanged(chat entities.Chat) {
	publisher.sendChatEvent(chat, "chat_changed")
}

type ChatEventsLoggingPublisher struct {
	publisher ports.ChatEventsPublisher
}

func (publisher ChatEventsLoggingPublisher) SendChatCreated(chat entities.Chat) {
	log.Printf("sending chat created event: %+v", chat)
	publisher.publisher.SendChatCreated(chat)
}

func (publisher ChatEventsLoggingPublisher) SendChatDeleted(chat entities.Chat) {
	log.Printf("sending chat deleted event: %+v", chat)
	publisher.publisher.SendChatDeleted(chat)
}

func (publisher ChatEventsLoggingPublisher) SendChatUserAction(chat entities.Chat) {
	log.Printf("sending chat user action event: %+v", chat)
	publisher.publisher.SendChatUserAction(chat)
}

func (publisher ChatEventsLoggingPublisher) SendChatChanged(chat entities.Chat) {
	log.Printf("sending chat changed event: %+v", chat)
	publisher.publisher.SendChatChanged(chat)
}

func NewChatEventsPublisher(connection rabbit.RabbitConnection) ports.ChatEventsPublisher {
	return ChatEventsLoggingPublisher{publisher: ChatEventsPublisher{connection: connection}}
}
