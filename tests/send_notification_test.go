package tests

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/google/uuid"
	"github.com/imhasandl/notification-service/cmd/server"
	"github.com/imhasandl/notification-service/internal/mocks"
	pb "github.com/imhasandl/notification-service/protos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendNotification(t *testing.T) {
	// Common test data
	validUserID := uuid.New()
	validDeviceToken := "device-token-123"

	// Create a notification payload
	notificationPayload := struct {
		Title          string    `json:"title"`
		SenderUsername string    `json:"sender_username"`
		ReceiverID     string    `json:"receiver_id"` // Fixed: ReceiverId -> ReceiverID
		Content        string    `json:"content"`
		SentAt         time.Time `json:"sent_at"`
	}{
		Title:          "Test Notification",
		SenderUsername: "testuser",
		ReceiverID:     validUserID.String(),
		Content:        "This is a test notification",
		SentAt:         time.Now(),
	}

	validNotificationBytes, _ := json.Marshal(notificationPayload)

	// Create invalid notification with bad UUID
	invalidUUIDNotification := notificationPayload
	invalidUUIDNotification.ReceiverID = "not-a-valid-uuid"
	invalidUUIDBytes, _ := json.Marshal(invalidUUIDNotification)

	// Define test cases
	testCases := []struct {
		name               string
		notificationBytes  []byte
		setupMocks         func(*mocks.MockQueries, *mocks.MockFirebaseClient)
		expectedErrMessage string
		shouldReturnError  bool
	}{
		{
			name:              "Success case",
			notificationBytes: validNotificationBytes,
			setupMocks: func(db *mocks.MockQueries, firebase *mocks.MockFirebaseClient) {
				db.On("GetDeviceTokensByUserID", mock.Anything, validUserID).
					Return(validDeviceToken, nil).Once()

				// Use FCMClient instead of directly using firebase
				firebase.FCMClient.On("Send",
					mock.Anything, // Match any context
					mock.MatchedBy(func(msg *messaging.Message) bool {
						return msg.Token == validDeviceToken &&
							msg.Notification.Title == notificationPayload.Title &&
							msg.Notification.Body == notificationPayload.Content
					}),
				).Return("message-id", nil).Once()
			},
			shouldReturnError: false,
		},
		{
			name:              "Invalid JSON",
			notificationBytes: []byte("invalid json"),
			setupMocks:        func(*mocks.MockQueries, *mocks.MockFirebaseClient) {},
			shouldReturnError: true,
		},
		{
			name:              "Invalid UUID",
			notificationBytes: invalidUUIDBytes,
			setupMocks:        func(*mocks.MockQueries, *mocks.MockFirebaseClient) {},
			shouldReturnError: true,
		},
		{
			name:              "Database error",
			notificationBytes: validNotificationBytes,
			setupMocks: func(db *mocks.MockQueries, firebase *mocks.MockFirebaseClient) {
				db.On("GetDeviceTokensByUserID", mock.Anything, validUserID).
					Return("", errors.New("database error")).Once()
			},
			shouldReturnError: true,
		},
		{
			name:              "Firebase error - should still succeed",
			notificationBytes: validNotificationBytes,
			setupMocks: func(db *mocks.MockQueries, firebase *mocks.MockFirebaseClient) {
				db.On("GetDeviceTokensByUserID", mock.Anything, validUserID).
					Return(validDeviceToken, nil).Once()

				// Use FCMClient instead of directly using firebase
				firebase.FCMClient.On("Send", mock.Anything, mock.Anything).
					Return("", errors.New("firebase error")).Once()
			},
			shouldReturnError: false, // Firebase errors are logged but don't fail the request
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh mocks for each test case
			mockDB := mocks.NewMockQueries()
			mockRabbitmq := mocks.NewMockRabbitMQ()
			mockFirebase := mocks.NewMockFirebaseClient()

			// Setup mocks according to test case
			tc.setupMocks(mockDB, mockFirebase)

			// Create server with mocks
			srv := server.NewServer(mockDB, mockRabbitmq, "test/path", mockFirebase)

			// Make the request
			request := &pb.SendNotificationRequest{
				Notification: tc.notificationBytes,
			}

			response, err := srv.SendNotification(context.Background(), request)

			// Check error expectation
			if tc.shouldReturnError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.True(t, response.Status)
			}

			// Verify all expectations were met
			mockDB.AssertExpectations(t)
			mockFirebase.AssertExpectations(t)
		})
	}
}
