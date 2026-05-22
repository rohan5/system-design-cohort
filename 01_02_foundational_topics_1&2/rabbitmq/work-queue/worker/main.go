package main

import (
	"fmt"
	"log"
	"strings"
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
	failOnError(err, "Failed on Dial")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed on Channel creation")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"work_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed on Queue declare")

	// err = ch.Qos(
	// 	1,
	// 	0,
	// 	false,
	// )
	// failOnError(err, "Failed on Qos setup")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed on consumer register")

	fmt.Println("Worker waiting for messages...")

	forever := make(chan struct{})

	go func() {
		for msg := range msgs {
			body := string(msg.Body)

			fmt.Printf("Received %s\n", body)

			// simulate work
			dots := strings.Count(body, ".")

			time.Sleep(time.Duration(dots) * time.Second)

			fmt.Println("Done")

			msg.Ack(false)
		}
	}()

	<-forever
}
