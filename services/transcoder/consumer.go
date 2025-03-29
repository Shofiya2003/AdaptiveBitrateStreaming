package main

import (
	"abr_backend/config"
	"abr_backend/utils"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type TranscodeJob struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

func main() {
	// Initialize RabbitMQ connection
	config.InitRabbitMq()

	// Start consuming messages
	StartConsumer()
}

// StartConsumer starts the queue consumer
func StartConsumer() {
	ch := config.Channel
	q := config.Queue

	msgs, err := ch.Consume(
		q.Name,                       // queue
		"transcode-service-consumer", // consumer
		false,                        // auto-ack
		false,                        // exclusive
		false,                        // no-local
		false,                        // no-wait
		nil,                          // args
	)
	if err != nil {
		log.Panicln("failed to register a consumer ", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Create a copy of the delivery for the goroutine
			delivery := d

			// Process message in a separate goroutine
			go func(msg amqp091.Delivery) {
				if err := processMessage(msg.Body); err != nil {
					log.Printf("Error processing message: %v", err)
					// Reject the message and requeue it
					msg.Nack(false, true) // false = don't requeue all messages, true = requeue this message
					return
				}
				// Acknowledge successful processing
				msg.Ack(false) // false = don't acknowledge all messages
			}(delivery)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// processMessage handles a single message from the queue
func processMessage(message []byte) error {
	var job TranscodeJob
	if err := json.Unmarshal(message, &job); err != nil {
		return fmt.Errorf("error unmarshaling message: %v", err)
	}

	// Run the transcoding process
	if err := utils.TranscodeVideo(job.Bucket, job.Key); err != nil {
		log.Printf("Error transcoding video: %v", err)
		return err
	}

	return nil
}
