package rabbit

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/streadway/amqp"
)

type queue struct {
	url          string
	name         string
	exchangeName string
	errorChannel chan *amqp.Error
	connection   *amqp.Connection
	channel      *amqp.Channel
	closed       bool
	consumers    []messageConsumer
}

type messageConsumer func([]byte)

func NewQueue(url string, qName string, exchangeName string) *queue {
	q := new(queue)
	q.url = url
	q.name = qName
	q.exchangeName = exchangeName
	q.consumers = make([]messageConsumer, 0)

	q.connect()
	go q.reconnector()

	return q
}

func (q *queue) Send(message string, contentType string) {
	err := q.channel.Publish(
		"",
		q.name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        []byte(message),
		},
	)

	logError("Sending message to queue failed", err)
}

func (q *queue) Consume(consumer messageConsumer) {
	log.Printf("Registering consumer...")
	deliveries, err := q.registerQueueConsumer()
	log.Printf("Consumer registered! Processing messages...")
	q.executeMessageConsumer(err, consumer, deliveries, false)
}

func (q *queue) Close() {
	log.Printf("Closing connection")
	q.closed = true
	q.channel.Close()
	q.connection.Close()
}

func (q *queue) reconnector() {
	for {
		err := <-q.errorChannel
		if !q.closed {
			logError("Reconnecting after connection closed", err)

			q.connect()
			q.recoverConsumers()
		}
	}
}

func (q *queue) connect() {
	for {
		log.Printf("Connecting to rabbitmq on %s", q.url)
		conn, err := amqp.Dial(q.url)
		if err == nil {
			q.connection = conn
			q.errorChannel = make(chan *amqp.Error)
			q.connection.NotifyClose(q.errorChannel)

			log.Println("Connection established!")

			q.openChannel()
			q.declareQueue()
			q.declareExchange()
			q.bindQueue()

			return
		}

		logError("Connection to rabbitmq failed. Retrying in 1 sec...", err)
		time.Sleep(1000 * time.Millisecond)
	}
}

func (q *queue) declareQueue() {
	_, err := q.channel.QueueDeclare(
		q.name,
		false,
		false,
		false,
		false,
		nil,
	)
	logError("Queue declaration failed", err)
}

func (q *queue) declareExchange() {
	err := q.channel.ExchangeDeclare(
		q.exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	logError("Exchange declaration failed", err)
}

func (q *queue) bindQueue() {
	err := q.channel.QueueBind(q.name, "", q.exchangeName, false, nil)
	logError("Exchange queue binding error", err)
}

func (q *queue) openChannel() {
	channel, err := q.connection.Channel()
	logError("Opening channel failed", err)
	q.channel = channel
}

func (q *queue) registerQueueConsumer() (<-chan amqp.Delivery, error) {
	msgs, err := q.channel.Consume(
		q.name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	logError("Consuming messages from queue failed", err)
	return msgs, err
}

func (q *queue) executeMessageConsumer(err error, consumer messageConsumer, deliveries <-chan amqp.Delivery, isRecovery bool) {
	if err != nil {
		return
	}

	if !isRecovery {
		q.consumers = append(q.consumers, consumer)
	}

	go func() {
		for delivery := range deliveries {
			consumer(delivery.Body[:])
		}
	}()
}

func (q *queue) recoverConsumers() {
	for i := range q.consumers {
		var consumer = q.consumers[i]

		log.Println("Recovering consumer...")
		msgs, err := q.registerQueueConsumer()
		log.Println("Consumer recovered! Continuing message processing...")
		q.executeMessageConsumer(err, consumer, msgs, true)
	}
}

func logError(message string, err error) {
	if err != nil {
		log.Printf("%s: %s", message, err)
		sentry.CaptureException(err)
	}
}
