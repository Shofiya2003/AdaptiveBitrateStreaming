package config

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Queue *amqp.Queue
var Channel *amqp.Channel

func InitRabbitMq() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	ch, err := conn.Channel()
	Channel = ch
	FailOnError(err, "Failed to open a channel")
	// defer ch.Close()

	queue, err := ch.QueueDeclare(
		"video_processing", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	Queue = &queue
	FailOnError(err, "Failed to declare a queue")

	log.Println("Successfully created a channel and a queue")

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
