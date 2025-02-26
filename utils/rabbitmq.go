package utils

import (
	"abr_backend/config"
	"abr_backend/data"
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishEvent(ch *amqp.Channel, q *amqp.Queue, body []byte) {
	if ch == nil || q == nil {
		log.Fatalln("Channel or Queue uninitialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	config.FailOnError(err, "Failed to publish a message")
	var event data.VideoEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Failed to decode message: %v", err)
		return
	}

	log.Printf(" [x] Sent %s\n", event.VideoURL)
}
