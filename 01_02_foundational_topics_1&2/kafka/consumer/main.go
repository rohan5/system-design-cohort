package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:29092"},
		GroupID:  "order-processors",
		Topic:    "orders",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	defer reader.Close()

	fmt.Println("Consumer started...")

	for {
		msg, err := reader.FetchMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"Consumed partition = %d offset = %d key = %s value =%s \n",
			msg.Partition,
			msg.Offset,
			msg.Key,
			msg.Value,
		)

		reader.CommitMessages(context.Background(), msg)
	}
}
