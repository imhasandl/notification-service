package server

import (
	"context"
	"log"

	"github.com/imhasandl/notification-service/internal/rabbitmq"
	pb "github.com/imhasandl/notification-service/protos"
)

// Consume starts consuming messages from the RabbitMQ queue and processes them as notifications.
// It registers a consumer with the RabbitMQ channel and handles incoming messages
// by converting them to notification requests and sending them through the SendNotification method.
// The method runs a goroutine that continuously processes messages from the queue.
func (s *Server) Consume() {
	msgs, err := s.rabbitmq.GetChannel().Consume(
		rabbitmq.QueueName, // queue
		"",                 // consumer
		false,              // auto-ack
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
			log.Printf("Received a message: %v", string(msg.Body))

			notificationReq := &pb.SendNotificationRequest{
				Notification: msg.Body,
			}

			_, err := s.SendNotification(context.Background(), notificationReq)
			if err != nil {
				log.Printf("Failed to send notification: %v", err)
				if rejectErr := msg.Reject(true); rejectErr != nil {
					log.Printf("Failed to reject message: %v", rejectErr)
				}
			}
		}
	}()
}
