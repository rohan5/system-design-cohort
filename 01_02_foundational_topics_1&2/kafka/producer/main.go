package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:29092"),
		Topic:    "orders",
		Balancer: &kafka.Hash{},
	}

	defer writer.Close()

	for i := 0; i < 20; i++ {
		// userId := "user-1"
		userId := fmt.Sprintf("user-%d", i%3)

		msg := kafka.Message{
			Key:   []byte(userId),
			Value: []byte(fmt.Sprintf("order-%d", i)),
		}

		err := writer.WriteMessages(context.Background(), msg)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"Produced message: key =%s value = %s\n", userId, msg.Value,
		)

		time.Sleep(1 * time.Second)
	}
}
