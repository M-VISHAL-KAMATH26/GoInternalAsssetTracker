package main

import (
	"log"
	"os"

	"asset-backend/internal/shared/mq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("test-consumer: failed to connect: %v", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("test-consumer: failed to open channel: %v", err)
	}
	defer channel.Close()

	msgs, err := channel.Consume(
		mq.ApprovalDecidedQueue,
		"",    // consumer tag
		true,  // auto-ack — fine for a throwaway test tool
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("test-consumer: failed to register consumer: %v", err)
	}

	log.Println("test-consumer: waiting for messages...")
	for msg := range msgs {
		log.Printf("received message: %s", msg.Body)
	}
}