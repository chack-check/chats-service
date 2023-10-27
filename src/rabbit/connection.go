package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
    if err != nil {
        log.Panicf("%s: %s", msg, err)
    }
}

type MessageEvent struct {
    Type string `json:"type"`
    MessageId int `json:"messageId"`
    IncludedUsers []int `json:"includedUsers"`
    ChatID int `json:"chatId"`
    SenderID int `json:"senderId"`
    MessageType string `json:"messageType"`
    Content string `json:"content"`
    VoiceURL string `json:"voiceUrl"`
    CircleURL string `json:"circleUrl"`
    Attachments []string `json:"attachments"`
    ReplyToID int `json:"replyToId"`
    Mentioned []int `json:"mentioned"`
    ReadedBy []int `json:"readedBy"`
}

type RabbitConnection struct {
    User string
    Pass string
    Host string
    Port string
    Connection *amqp.Connection
    Channel *amqp.Channel
    Queue amqp.Queue
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

func (conn *RabbitConnection) DeclareQueue(queueName string) {
    queue, err := conn.Channel.QueueDeclare(
        queueName,
        false,
        false,
        false,
        false,
        nil,
    )
    failOnError(err, "Failed to declare a queue")
    conn.Queue = queue
}

func (conn *RabbitConnection) SendEvent(event []byte) error {
    if conn.Connection.IsClosed() {
        conn.Connect()
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    log.Printf("Sending content to rabbitmq: %s, queue name: %s, closed: %v", event, conn.Queue.Name, conn.Connection.IsClosed())
    return conn.Channel.PublishWithContext(
        ctx,
        "",
        conn.Queue.Name,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body: event,
        },
    )
}

func (conn *RabbitConnection) SendMessageEvent(message *MessageEvent) error {
    body, err := json.Marshal(message)
    if err != nil {
        return err
    }
    fmt.Printf("Sending message to queue: %v\n", message)
    return conn.SendEvent(body)
}

func (conn *RabbitConnection) Close() {
    conn.Connection.Close()
    conn.Channel.Close()
}

func NewEventsRabbitConnection() *RabbitConnection {
    user := os.Getenv("RABBIT_USER")
    pass := os.Getenv("RABBIT_PASS")
    host := os.Getenv("RABBIT_HOST")
    port := os.Getenv("RABBIT_PORT")
    queueName := os.Getenv("RABBIT_EVENTS_QUEUE_NAME")
    conn := &RabbitConnection{
        User: user,
        Pass: pass,
        Host: host,
        Port: port,
    }
    conn.Connect()
    conn.DeclareQueue(queueName)
    return conn
}

var EventsRabbitConnection *RabbitConnection = NewEventsRabbitConnection()
