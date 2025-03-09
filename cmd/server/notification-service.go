package server

import (
	"context"
	"log"

	"github.com/imhasandl/notification-service/internal/database"
	"github.com/imhasandl/notification-service/internal/rabbitmq"
	pb "github.com/imhasandl/notification-service/protos"
)

type server struct {
	pb.UnimplementedNotificationServiceServer
	db          *database.Queries
	tokenSecret string
	rabbitmq    *rabbitmq.RabbitMQ
}

func NewServer(db *database.Queries, tokenSecret string, rabbitmq *rabbitmq.RabbitMQ) *server {
	return &server{
		pb.UnimplementedNotificationServiceServer{},
		db,
		tokenSecret,
		rabbitmq,
	}
}

func (s *server) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error){
	log.Printf("Sending notification with message: %s", req.GetNotification())

	return &pb.SendNotificationResponse{
		Success: true,
	}, nil
}