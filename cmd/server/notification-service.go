package server

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/imhasandl/notification-service/cmd/helper"
	"github.com/imhasandl/notification-service/internal/database"
	"github.com/imhasandl/notification-service/internal/rabbitmq"
	pb "github.com/imhasandl/notification-service/protos"
	"google.golang.org/grpc/codes"
)

type server struct {
	pb.UnimplementedNotificationServiceServer
	db          *database.Queries
	tokenSecret string
	rabbitmq    *rabbitmq.RabbitMQ
}

type Notification struct {
	ReceiverId string    `json:"receiver_id"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
}

func NewServer(db *database.Queries, tokenSecret string, rabbitmq *rabbitmq.RabbitMQ) *server {
	return &server{
		pb.UnimplementedNotificationServiceServer{},
		db,
		tokenSecret,
		rabbitmq,
	}
}

func (s *server) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	log.Printf("Sending notification with message: %s", req.GetNotification())

	var notification Notification
	err := json.Unmarshal(req.GetNotification(), &notification)
	if err != nil {
		 return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't parse json - SendNotification", err)
	}

	return &pb.SendNotificationResponse{
		Success: true,
	}, nil
}
