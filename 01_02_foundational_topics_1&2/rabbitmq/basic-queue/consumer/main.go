package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare queue")

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer name
		false, // auto-ack
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register consumer")

	fmt.Println("Waiting for messages...")

	forever := make(chan struct{})

	go func() {
		for msg := range msgs {
			fmt.Printf("Received message: %s\n", msg.Body)

			// simulate work
			time.Sleep(2 * time.Second)

			fmt.Println("Processing done")

			msg.Ack(false)
		}
	}()

	<-forever
}
