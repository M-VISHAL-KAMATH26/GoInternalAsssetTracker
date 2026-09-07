package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer wraps a RabbitMQ connection/channel for consuming messages,
// mirroring Publisher's shape. Unlike the throwaway test-consumer's
// auto-ack, this uses manual ack/nack so a failed handler can requeue
// the message instead of silently losing it.
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConsumer(rabbitURL string) (*Consumer, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Process one message at a time before requiring an ack — avoids
	// RabbitMQ flooding this consumer with every queued message at once.
	if err := channel.Qos(1, 0, false); err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return &Consumer{conn: conn, channel: channel}, nil
}

// Consume registers handler to be called for every message on the
// given queue. handler returning nil acks the message; returning an
// error nacks it with requeue=true, so RabbitMQ redelivers it later.
// This call blocks the calling goroutine until the channel/connection closes.
func (c *Consumer) Consume(queueName string, handler func(body []byte) error) error {
	_, err := c.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := c.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	for msg := range msgs {
		if err := handler(msg.Body); err != nil {
			msg.Nack(false, true) // requeue on failure
			continue
		}
		msg.Ack(false)
	}

	return nil
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}