package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

const (
	RoutingKey   = "message"
	QueueName    = "message_queue"
	ExchangeName = "message_exchange"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

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

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

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
