package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	body := "Hello from work-queue"

	if len(os.Args) > 1 {
		body = strings.Join(os.Args[1:], " ")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	failOnError(err, "Msg piblish failed")

	fmt.Printf("Sent: %s\n", body)

}
