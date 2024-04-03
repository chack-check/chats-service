package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/chack-check/chats-service/settings"
	"github.com/getsentry/sentry-go"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartConsumer(ctag string) error {
	for {
		dialString := fmt.Sprintf(
			"amqp://%s:%s@%s:%d/",
			settings.Settings.APP_RABBIT_USER,
			settings.Settings.APP_RABBIT_PASSWORD,
			settings.Settings.APP_RABBIT_HOST,
			settings.Settings.APP_RABBIT_PORT,
		)
		log.Printf("dialing %q", dialString)
		conn, err := amqp.Dial(dialString)
		if err != nil {
			log.Printf("Error dialing rabbitmq conenction: %v", err)
			sentry.CaptureException(err)
			return fmt.Errorf("dial: %s", err)
		}

		go func() {
			fmt.Printf("closing: %s", <-conn.NotifyClose(make(chan *amqp.Error)))
		}()

		log.Printf("got Connection, getting Channel")
		channel, err := conn.Channel()
		if err != nil {
			log.Printf("Error getting channel from connection: %v", err)
			sentry.CaptureException(err)
			return fmt.Errorf("channel: %s", err)
		}

		channel.Qos(10, 0, false)

		log.Printf("got Channel, declaring Exchange %q", settings.Settings.APP_RABBIT_USERS_EXCHANGE_NAME)
		err = channel.ExchangeDeclare(
			settings.Settings.APP_RABBIT_USERS_EXCHANGE_NAME,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("Error declaring exchange %s: %v", settings.Settings.APP_RABBIT_USERS_EXCHANGE_NAME, err)
			sentry.CaptureException(err)
			return fmt.Errorf("exchange: %s", err)
		}

		log.Printf("got Exchange, declaring Queue %q", settings.Settings.APP_RABBIT_CONSUMER_QUEUE_NAME)
		queue, err := channel.QueueDeclare(
			settings.Settings.APP_RABBIT_CONSUMER_QUEUE_NAME, // name of the queue
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // noWait
			nil,   // arguments
		)
		if err != nil {
			log.Printf("Error declaring queue %s: %v", settings.Settings.APP_RABBIT_CONSUMER_QUEUE_NAME, err)
			sentry.CaptureException(err)
			return fmt.Errorf("queue Declare: %s", err)
		}

		log.Printf("declared Queue (%q %d messages, %d consumers). Binding queue to exchange", queue.Name, queue.Messages, queue.Consumers)

		err = channel.QueueBind(
			settings.Settings.APP_RABBIT_CONSUMER_QUEUE_NAME,
			"",
			settings.Settings.APP_RABBIT_USERS_EXCHANGE_NAME,
			false,
			nil,
		)
		if err != nil {
			log.Printf("Error binding queue for chats exchange: %v", err)
			sentry.CaptureException(err)
			return fmt.Errorf("queue Bind: %s", err)
		}

		log.Printf("starting Consume (consumer tag %q)", ctag)
		deliveries, err := channel.Consume(
			queue.Name, // name
			ctag,       // consumerTag,
			false,      // noAck
			false,      // exclusive
			false,      // noLocal
			false,      // noWait
			nil,        // arguments
		)
		if err != nil {
			log.Printf("Error starting consume: %v", err)
			sentry.CaptureException(err)
			return fmt.Errorf("queue Consume: %s", err)
		}

		for d := range deliveries {
			log.Printf("Received message: %q", d.Body)
			var message_type DataStruct
			err := json.Unmarshal(d.Body, &message_type)
			if err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				sentry.CaptureException(err)
				d.Ack(true)
				continue
			}

			if message_type.EventType == "user_created" {
				HandleUserCreated(d.Body)
			} else {
				log.Printf("Received unhandled rabbitmq event: %q", d.Body)
			}

			d.Ack(true)
		}
		sentry.CaptureMessage("Restarting rabbitmq connection")
		log.Printf("handle: deliveries channel closed")
	}
}
