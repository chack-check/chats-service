package rabbit

import (
	"encoding/json"
	"log"

	"chats-service/configs"
	"chats-service/internal/application/ports"

	"github.com/getsentry/sentry-go"
)

func StartConsumer(ctag string, messageEventsPublisher ports.MessageEventsPublisher) error {
	configuration := configs.GetRabbitConfiguration()
	queue := NewQueue(configuration.Host, configuration.ConsumerQueueName, configuration.UsersExchangeName)
	recognitionQueue := NewQueue(configuration.Host, configuration.RecognitionQueueName, configuration.RecognitionExchangeName)

	queue.Consume(func(msg []byte) {
		log.Printf("fetched event: %s", string(msg))
		var event SystemEvent
		err := json.Unmarshal(msg, &event)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			sentry.CaptureException(err)
			return
		}

		if event.EventType == "user_created" {
			log.Printf("Fetched user created event: %+v", event)
			HandleUserCreated(event)
		}
	})

	recognitionQueue.Consume(func(msg []byte) {
		log.Printf("fetched recognition event: %s", string(msg))
		var event RecognitionEvent
		err := json.Unmarshal(msg, &event)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			sentry.CaptureException(err)
			return
		}

		HandleMessageRecognized(event.MessageID, event.Content, messageEventsPublisher)
	})

	return nil
}
