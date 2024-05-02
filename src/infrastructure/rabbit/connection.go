package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type EventSavedFile struct {
	OriginalUrl       string  `json:"originalUrl"`
	OriginalFilename  string  `json:"originalFilename"`
	ConvertedUrl      *string `json:"convertedUrl"`
	ConvertedFilename *string `json:"convertedFilename"`
}

type EventActionUser struct {
	Id         int     `json:"id"`
	LastName   string  `json:"lastName"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	Username   string  `json:"username"`
}

type ChatEvent struct {
	Id         int                          `json:"id"`
	Avatar     *EventSavedFile              `json:"avatar"`
	Title      string                       `json:"title"`
	Type       string                       `json:"type"`
	Members    []int                        `json:"members"`
	IsArchived bool                         `json:"isArchived"`
	OwnerId    int                          `json:"ownerId"`
	Admins     []int                        `json:"admins"`
	Actions    map[string][]EventActionUser `json:"actions"`
}

type EventMessageReaction struct {
	UserId  int    `json:"userId"`
	Content string `json:"content"`
}

type MessageEvent struct {
	Id          int                    `json:"id"`
	SenderId    int                    `json:"senderId"`
	ChatId      int                    `json:"chatId"`
	Type        string                 `json:"type"`
	Content     *string                `json:"content"`
	Voice       *EventSavedFile        `json:"voice"`
	Circle      *EventSavedFile        `json:"circle"`
	Attachments []EventSavedFile       `json:"attachments"`
	ReplyToId   *int                   `json:"replyToId"`
	Mentioned   []int                  `json:"mentioned"`
	ReadedBy    []int                  `json:"readedBy"`
	Reactions   []EventMessageReaction `json:"reactions"`
	CreatedAt   *time.Time             `json:"createdAt"`
}

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

func (conn *RabbitConnection) Close() {
	conn.Connection.Close()
	conn.Channel.Close()
}

func NewEventsRabbitConnection(user string, pass string, host string, port int, exchangeName string) *RabbitConnection {
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

var EventsRabbitConnection *RabbitConnection = NewEventsRabbitConnection(
	Settings.APP_RABBIT_USER,
	Settings.APP_RABBIT_PASSWORD,
	Settings.APP_RABBIT_HOST,
	Settings.APP_RABBIT_PORT,
	Settings.APP_RABBIT_PUBLISHER_EXCHANGE_NAME,
)
