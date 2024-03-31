package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chack-check/chats-service/settings"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type MessageEvent struct {
	Type          string   `json:"type"`
	MessageId     int      `json:"message_id"`
	IncludedUsers []int    `json:"included_users"`
	ChatID        int      `json:"chat_id"`
	SenderID      int      `json:"sender_id"`
	MessageType   string   `json:"message_type"`
	Content       string   `json:"content"`
	VoiceURL      string   `json:"voice_url"`
	CircleURL     string   `json:"circle_url"`
	Attachments   []string `json:"attachments"`
	ReplyToID     int      `json:"reply_to_id"`
	Mentioned     []int    `json:"mentioned"`
	ReadedBy      []int    `json:"readed_by"`
	Datetime      string   `json:"datetime"`
}

type ChatEvent struct {
	Type          string `json:"type"`
	ChatId        int    `json:"chatId"`
	IncludedUsers []int  `json:"includedUsers"`
	AvatarURL     string `json:"avatarURL"`
	Title         string `json:"title"`
	ChatType      string `json:"chatType"`
	Members       []int  `json:"members"`
	IsArchived    bool   `json:"isArchived"`
	OwnerID       int    `json:"ownerId"`
	Admins        []int  `json:"admins"`
}

type ReadMessageEvent struct {
	Type          string `json:"type"`
	MessageId     int    `json:"messageId"`
	ReadedBy      []int  `json:"readedBy"`
	ChatID        int    `json:"chatId"`
	IncludedUsers []int  `json:"includedUsers"`
}

type IRabbitConnection interface {
	Connect()
	DeclareExchange()
	SendEvent(event interface{}) error
	Close()
}

type MockRabbitConnection struct{}

type RabbitConnection struct {
	User         string
	Pass         string
	Host         string
	Port         string
	ExchangeName string
	Connection   *amqp.Connection
	Channel      *amqp.Channel
}

func (conn *RabbitConnection) Connect() {
	dialString := fmt.Sprintf("amqp://%s:%s@%s:%s/", conn.User, conn.Pass, conn.Host, conn.Port)
	connection, err := amqp.Dial(dialString)
	failOnError(err, "Failed to connect to RabbitMQ")
	conn.Connection = connection

	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	conn.Channel = channel
}

func (conn *MockRabbitConnection) Connect() {
	log.Print("Connect to rabbitmq")
}

func (conn *RabbitConnection) DeclareExchange() {
	err := conn.Channel.ExchangeDeclare(
		conn.ExchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")
}

func (conn *MockRabbitConnection) DeclareExchange() {
	log.Printf("Declaring exchange %s", "test")
}

func (conn *RabbitConnection) SendEvent(event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	log.Printf("Sending message to queue: %+v", event)

	if conn.Connection.IsClosed() {
		log.Printf("Rabbitmq connection is closed. Reconnecting")
		conn.Connect()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Printf("Sending content to rabbitmq: %s, exchange name: %s, closed: %v", body, conn.ExchangeName, conn.Connection.IsClosed())
	return conn.Channel.PublishWithContext(
		ctx,
		conn.ExchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (conn *MockRabbitConnection) SendEvent(event interface{}) error {
	log.Printf("Sending event to queue: %v", event)
	return nil
}

func (conn *RabbitConnection) Close() {
	conn.Connection.Close()
	conn.Channel.Close()
}

func (conn *MockRabbitConnection) Close() {
	log.Print("Closed rabbitmq connection")
}

func NewEventsRabbitConnection(user string, pass string, host string, port int, exchangeName string) IRabbitConnection {
	log.Printf("Initializing rabbitmq connection for environment: %s", settings.Settings.APP_ENVIRONMENT)
	if settings.Settings.APP_ENVIRONMENT == "test" {
		conn := &MockRabbitConnection{}
		return conn
	}

	conn := &RabbitConnection{
		User:         user,
		Pass:         pass,
		Host:         host,
		Port:         fmt.Sprint(port),
		ExchangeName: exchangeName,
	}
	conn.Connect()
	conn.DeclareExchange()
	log.Printf("Declared rabbitmq connection: %+v", conn)
	return conn
}

var EventsRabbitConnection IRabbitConnection = NewEventsRabbitConnection(
	settings.Settings.APP_RABBIT_USER,
	settings.Settings.APP_RABBIT_PASSWORD,
	settings.Settings.APP_RABBIT_HOST,
	settings.Settings.APP_RABBIT_PORT,
	settings.Settings.APP_RABBIT_PUBLISHER_EXCHANGE_NAME,
)
