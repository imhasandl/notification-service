package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	App *firebase.App
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
		App: app,
		FCMClient: fcmClient,
	}, nil
}