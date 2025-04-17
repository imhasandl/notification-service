package tests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/imhasandl/notification-service/cmd/server"
	"github.com/imhasandl/notification-service/internal/mocks"
	pb "github.com/imhasandl/notification-service/protos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSendNotification(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name          string
		receiverID    string
		deviceToken   string
		expectError   bool
		errorContains string
		setupMocks    func(*mocks.MockDBQuerier, *mocks.MockFCMClient)
	}{
		{
			name:        "Success case",
			receiverID:  "f6b3f9cf-7e9c-48fe-aa1c-e5afbef59770",
			deviceToken: "device-token-123",
			expectError: false,
			setupMocks: func(db *mocks.MockDBQuerier, fcm *mocks.MockFCMClient) {
				receiverID, _ := uuid.Parse("f6b3f9cf-7e9c-48fe-aa1c-e5afbef59770")
				db.On("GetDeviceTokensByUserID", mock.Anything, receiverID).Return("device-token-123", nil)

				// The key fix - use correct argument matchers
				fcm.On("Send", mock.Anything, mock.AnythingOfType("*messaging.Message")).Return("message-id", nil)
			},
		},
		{
			name:          "Invalid UUID",
			receiverID:    "invalid-uuid",
			expectError:   true,
			errorContains: "can't parse receiver id",
			setupMocks:    func(*mocks.MockDBQuerier, *mocks.MockFCMClient) {},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockDB := new(mocks.MockDBQuerier)
			mockRabbitMQ := new(mocks.MockRabbitMQClient)
			mockFirebase := new(mocks.MockFirebaseClient)
			mockFCM := new(mocks.MockFCMClient)

			// Setup mockFirebase to return mockFCM
			mockFirebase.On("GetMessagingClient").Return(mockFCM)

			// Setup test-specific mocks
			tc.setupMocks(mockDB, mockFCM)

			// Create server with mocks
			srv := server.NewServer(mockDB, mockRabbitMQ, "test-path", mockFirebase)

			// Create notification payload
			notification := server.Notification{
				Title:          "Test Notification",
				SenderUsername: "testuser",
				ReceiverID:     tc.receiverID,
				Content:        "This is a test notification",
				SentAt:         time.Now(),
			}

			notificationBytes, err := json.Marshal(notification)
			require.NoError(t, err)

			// Create request
			req := &pb.SendNotificationRequest{
				Notification: notificationBytes,
			}

			// Call the method
			resp, err := srv.SendNotification(context.Background(), req)

			// Check results
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Status)
			}

			// Verify mocks
			mockDB.AssertExpectations(t)
			mockFirebase.AssertExpectations(t)
			mockFCM.AssertExpectations(t)
		})
	}
}
