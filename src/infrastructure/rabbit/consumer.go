package rabbit

import (
	"encoding/json"
	"log"

	"github.com/getsentry/sentry-go"
)

func StartConsumer(ctag string) error {
	queue := NewQueue(Settings.APP_RABBIT_HOST, Settings.APP_RABBIT_CONSUMER_QUEUE_NAME, Settings.APP_RABBIT_USERS_EXCHANGE_NAME)

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

	return nil
}
