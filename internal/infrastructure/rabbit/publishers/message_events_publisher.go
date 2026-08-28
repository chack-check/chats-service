package publishers

import (
	"log"

	"chats-service/internal/application/ports"
	"chats-service/internal/domain/entities"
	"chats-service/internal/infrastructure/rabbit"
)

type MessageEventsPublisher struct {
	connection rabbit.RabbitConnection
}

func (publisher MessageEventsPublisher) getSystemEventForMessage(message entities.Message, eventType string) (*rabbit.SystemEvent, error) {
	messageEvent := rabbit.MessageToMessageEvent(message)
	chat := message.GetChat()
	return rabbit.NewSystemEvent(eventType, chat.GetMembers(), messageEvent)
}

func (publisher MessageEventsPublisher) sendMessageEvent(message entities.Message, eventType string) {
	systemEvent, err := publisher.getSystemEventForMessage(message, eventType)
	if err != nil {
		return
	}

	publisher.connection.SendEvent(systemEvent)
}

func (publisher MessageEventsPublisher) SendMessageReacted(message entities.Message) {
	publisher.sendMessageEvent(message, "message_reacted")
}

func (publisher MessageEventsPublisher) SendReactionDeleted(message entities.Message) {
	publisher.sendMessageEvent(message, "message_reaction_deleted")
}

func (publisher MessageEventsPublisher) SendMessageReaded(message entities.Message) {
	publisher.sendMessageEvent(message, "message_readed")
}

func (publisher MessageEventsPublisher) SendMessageDeleted(message entities.Message) {
	publisher.sendMessageEvent(message, "message_deleted")
}

func (publisher MessageEventsPublisher) SendMessageUpdated(message entities.Message) {
	publisher.sendMessageEvent(message, "message_updated")
}

func (publisher MessageEventsPublisher) SendMessageCreated(message entities.Message) {
	publisher.sendMessageEvent(message, "message_created")
}

type MessageEventsLoggingPublisher struct {
	publisher ports.MessageEventsPublisher
}

func (publisher MessageEventsLoggingPublisher) SendMessageReacted(message entities.Message) {
	log.Printf("sending message reacted event: %+v", message)
	publisher.publisher.SendMessageReacted(message)
}

func (publisher MessageEventsLoggingPublisher) SendReactionDeleted(message entities.Message) {
	log.Printf("sending message reaction deleted event: %+v", message)
	publisher.publisher.SendReactionDeleted(message)
}

func (publisher MessageEventsLoggingPublisher) SendMessageReaded(message entities.Message) {
	log.Printf("sending message readed event: %+v", message)
	publisher.publisher.SendMessageReaded(message)
}

func (publisher MessageEventsLoggingPublisher) SendMessageDeleted(message entities.Message) {
	log.Printf("sending message deleted event: %+v", message)
	publisher.publisher.SendMessageDeleted(message)
}

func (publisher MessageEventsLoggingPublisher) SendMessageUpdated(message entities.Message) {
	log.Printf("sending message updated event: %+v", message)
	publisher.publisher.SendMessageUpdated(message)
}

func (publisher MessageEventsLoggingPublisher) SendMessageCreated(message entities.Message) {
	log.Printf("sending message created event: %+v", message)
	publisher.publisher.SendMessageCreated(message)
}

func NewMessageEventsPublisher(connection rabbit.RabbitConnection) ports.MessageEventsPublisher {
	return MessageEventsLoggingPublisher{publisher: MessageEventsPublisher{connection: connection}}
}
