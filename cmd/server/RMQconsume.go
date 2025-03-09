package server

import (
	"context"
	"log"

	"github.com/imhasandl/notification-service/internal/rabbitmq"
	pb "github.com/imhasandl/notification-service/protos"
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

			notificationReq := &pb.SendNotificationRequest{
				Notification: msg.Body,
			}

			_, err := s.SendNotification(context.Background(), notificationReq)
			if err != nil {
				log.Printf("Failed to send notification: %v", err)
			}
		}
	}()
}
