package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

const (
	// ExchangeName is the name of the RabbitMQ exchange used for notifications
	ExchangeName = "notifications.topic"
	// QueueName is the name of the queue for processing notifications
	QueueName = "notification_service_queue"
)

// RabbitMQ represents a RabbitMQ client connection
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// Client defines the interface for RabbitMQ operations
type Client interface {
	Close()
	GetChannel() *amqp.Channel
}

// Ensure RabbitMQ implements the interface
var _ Client = (*RabbitMQ)(nil)

// NewRabbitMQ creates a new RabbitMQ client connected to the specified URL
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("can't connect to rabbit mq: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("can't connect to the channel: %v", err)
		return nil, err
	}

	// Declare the topic exchange
	if err := ch.ExchangeDeclare(
		ExchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		log.Printf("failed to declare exchange: %v", err)
		return nil, err
	}

	// Declare the queue for the notification service
	queue, err := ch.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("failed to declare queue: %v", err)
		return nil, err
	}

	// Bind the queue to the exchange with a wildcard routing key
	if err := ch.QueueBind(
		queue.Name,   // queue name
		"#",          // routing key - match all messages
		ExchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		log.Printf("failed to bind queue: %v", err)
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			log.Printf("error closing channel: %v", err)
		}
	}
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			log.Printf("error closing connection: %v", err)
		}
	}
}

// GetChannel returns the current channel
func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.Channel
}
