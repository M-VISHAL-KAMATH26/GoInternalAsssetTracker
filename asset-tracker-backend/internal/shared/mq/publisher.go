package mq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const ApprovalDecidedQueue = "approval_decided"

// Publisher wraps a RabbitMQ connection/channel and exposes a simple
// Publish method. Any service that needs to publish events can use
// this same wrapper.
type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewPublisher connects to RabbitMQ and declares the queue this
// project uses, so it's guaranteed to exist before anyone publishes
// or consumes from it.
func NewPublisher(rabbitURL string) (*Publisher, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = channel.QueueDeclare(
		ApprovalDecidedQueue,
		true,  // durable — survives a RabbitMQ restart
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Publisher{conn: conn, channel: channel}, nil
}

func (p *Publisher) PublishApprovalDecided(ctx context.Context, event ApprovalDecidedEvent) error {
	body, err := event.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return p.channel.PublishWithContext(ctx,
		"",                  // exchange — "" means the default exchange
		ApprovalDecidedQueue, // routing key — with the default exchange, this is just the queue name
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) Close() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	return p.conn.Close()
}