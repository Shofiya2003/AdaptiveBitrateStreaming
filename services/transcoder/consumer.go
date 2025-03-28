package main

import (
	"abr_backend/config"
	"abr_backend/utils"
	"encoding/json"
	"fmt"
	"log"
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
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicln("failed to register a consumer ", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			go processMessage(d.Body)
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
