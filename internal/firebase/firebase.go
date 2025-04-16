package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// ClientInterface defines the interface for Firebase operations
type ClientInterface interface {
	// Add the methods needed by server
	GetMessagingClient() MessagingClient
}

// MessagingClient defines the interface for Firebase messaging operations
type MessagingClient interface {
	Send(ctx context.Context, message *messaging.Message) (string, error)
}

// Client represents a Firebase client with messaging capabilities
type Client struct {
	App       *firebase.App
	FCMClient *messaging.Client
}

// InitFirebase initializes a new Firebase client with the provided credentials
func InitFirebase(ctx context.Context, credentialsFilePath string) (*Client, error) {
	opt := option.WithCredentialsFile(credentialsFilePath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("Error initializing Firebase app: %v", err)
		return nil, err
	}

	// Init FCM Client
	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("Error getting Firebase Messaging client: %v", err)
		return nil, err
	}

	return &Client{
		App:       app,
		FCMClient: fcmClient,
	}, nil
}

// GetMessagingClient returns the Firebase Cloud Messaging client
func (fc *Client) GetMessagingClient() MessagingClient {
	return fc.FCMClient
}
