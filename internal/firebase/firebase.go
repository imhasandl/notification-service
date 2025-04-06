package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// FirebaseClientInterface defines the interface for Firebase operations
type FirebaseClientInterface interface {
	// Add the methods needed by server
	GetMessagingClient() MessagingClient
}

// MessagingClient defines the interface for Firebase messaging operations
type MessagingClient interface {
	Send(ctx context.Context, message *messaging.Message) (string, error)
}

type FirebaseClient struct {
	App       *firebase.App
	FCMClient *messaging.Client
}

func InitFirebase(ctx context.Context, credentialsFilePath string) (*FirebaseClient, error) {
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

	return &FirebaseClient{
		App:       app,
		FCMClient: fcmClient,
	}, nil
}

// Ensure FirebaseClient implements the interface
func (fc *FirebaseClient) GetMessagingClient() MessagingClient {
	return fc.FCMClient
}
