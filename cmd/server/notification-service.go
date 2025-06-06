package server

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"github.com/imhasandl/notification-service/cmd/helper"
	"github.com/imhasandl/notification-service/internal/database"
	"github.com/imhasandl/notification-service/internal/firebase"
	"github.com/imhasandl/notification-service/internal/rabbitmq"
	pb "github.com/imhasandl/notification-service/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DBQuerier defines the interface for database operations required by the notification service
type DBQuerier interface {
	GetDeviceTokensByUserID(ctx context.Context, userID uuid.UUID) (string, error)
	RegisterDeviceToken(ctx context.Context, arg database.RegisterDeviceTokenParams) (database.DeviceToken, error)
	DeleteDeviceToken(ctx context.Context, arg database.DeleteDeviceTokenParams) error
	SendNotification(ctx context.Context) error
}

// Server implements the notification service gRPC server
type Server struct {
	pb.UnimplementedNotificationServiceServer
	db              DBQuerier
	rabbitmq        rabbitmq.Client
	firebaseKeyPath string
	firebase        firebase.ClientInterface
}

// Notification represents the structure of a notification message
type Notification struct {
	Title          string    `json:"title"`
	SenderUsername string    `json:"sender_username"`
	ReceiverID     string    `json:"receiver_id"` // Fixed: ReceiverId -> ReceiverID
	Content        string    `json:"content"`
	SentAt         time.Time `json:"sent_at"`
}

// NewServer creates a new notification service server with the provided dependencies
func NewServer(db DBQuerier, rabbitmq rabbitmq.Client, firebaseKeyPath string, firebase firebase.ClientInterface) *Server {
	return &Server{
		pb.UnimplementedNotificationServiceServer{},
		db,
		rabbitmq,
		firebaseKeyPath,
		firebase,
	}
}

// SendNotification handles requests to send push notifications to users.
// It implements the NotificationServiceServer interface from the protobuf definition.
func (s *Server) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	var notification Notification
	err := json.Unmarshal(req.GetNotification(), &notification)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't parse json - SendNotification", err)
	}

	receiverID, err := uuid.Parse(notification.ReceiverID) // Fixed: receiverId -> receiverID
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't parse receiver id - SendNotification", err)
	}

	log.Printf("Sent push notification to: %v", notification.ReceiverID)

	// Check if Firebase isn't initialized
	if s.firebase == nil && s.firebase.GetMessagingClient() == nil {
		log.Printf("Firebase not initialized, skipping push notification")
	}

	receiverDeviceToken, err := s.db.GetDeviceTokensByUserID(ctx, receiverID)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "Can't get receiver's device token from db - SendNotification", err)
	}

	if receiverDeviceToken == "" {
		log.Printf("No device token found for user %s", notification.ReceiverID)
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: notification.Title,
			Body:  notification.Content,
		},
		Token: receiverDeviceToken,
		Data: map[string]string{
			"title":           notification.Title,
			"sender_username": notification.SenderUsername,
			"receiver_id":     notification.ReceiverID,
			"content":         notification.Content,
			"sent_at":         notification.SentAt.Format(time.RFC3339),
		},
	}

	// Send the Message
	response, err := s.firebase.GetMessagingClient().Send(ctx, message)
	if err != nil {
		log.Printf("Error sending FCM message: %v", err)
	} else {
		log.Printf("Successfully sent FCM message: %s", response)
	}

	return &pb.SendNotificationResponse{
		Status: true,
	}, nil
}

// RegisterDeviceToken handles requests to register a new device token for push notifications.
// It implements the NotificationServiceServer interface from the protobuf definition.
func (s *Server) RegisterDeviceToken(ctx context.Context, req *pb.RegisterDeviceTokenRequest) (*pb.RegisterDeviceTokenResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't parse provided uuid - RegisterDeviceToken", err)
	}

	deviceTokenParams := database.RegisterDeviceTokenParams{
		ID:          uuid.New(),
		UserID:      userID,
		DeviceToken: req.GetDeviceToken(),
		DeviceType:  req.GetDeviceType(),
	}

	deviceToken, err := s.db.RegisterDeviceToken(ctx, deviceTokenParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't get device token from db - RegisterDeviceToken", err)
	}

	return &pb.RegisterDeviceTokenResponse{
		DeviceToken: &pb.DeviceToken{
			Id:          deviceToken.ID.String(),
			UserId:      deviceToken.UserID.String(),
			DeviceToken: deviceToken.DeviceToken,
			DeviceType:  deviceToken.DeviceType,
			CreatedAt:   timestamppb.New(deviceToken.CreatedAt),
			UpdatedAt:   timestamppb.New(deviceToken.UpdatedAt),
		},
	}, nil
}

// DeleteDeviceToken handles requests to delete a device token for a user.
// It implements the NotificationServiceServer interface from the protobuf definition.
func (s *Server) DeleteDeviceToken(ctx context.Context, req *pb.DeleteDeviceTokenRequest) (*pb.DeleteDeviceTokenResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't parse user's incoming id - DeleteDeviceToken", err)
	}

	deleteDeviceTokenParams := database.DeleteDeviceTokenParams{
		UserID:      userID,
		DeviceToken: req.GetDeviceToken(),
	}

	err = s.db.DeleteDeviceToken(ctx, deleteDeviceTokenParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't delete device token of a user - DeleteDeviceToken", err)
	}

	return &pb.DeleteDeviceTokenResponse{
		Status: true,
	}, nil
}
