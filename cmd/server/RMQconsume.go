package server

import (
	"log"

	"github.com/imhasandl/notification-service/internal/rabbitmq"
)

func (s *server) Consume() {
	msgs, err := s.rabbitmq.Channel.Consume(
		rabbitmq.QueueName, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %v", msg.Body)
		}
	}()
}
