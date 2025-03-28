package utils

import (
	"abr_backend/config"
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishEvent(ch *amqp.Channel, q *amqp.Queue, body []byte) error {
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

	if err != nil {
		config.FailOnError(err, "Failed to publish a message")
		return err
	}

	return nil
}
