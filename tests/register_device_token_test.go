package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/imhasandl/notification-service/cmd/server"
	"github.com/imhasandl/notification-service/internal/mocks"
	"github.com/imhasandl/notification-service/internal/database"
	pb "github.com/imhasandl/notification-service/protos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterDeviceToken(t *testing.T) {
	// Common test data
	validUserID := uuid.New()
	validDeviceToken := "device-token-123"
	validDeviceType := "android"

	// Define test cases
	testCases := []struct {
		name              string
		userID            string // Fixed: userId -> userID
		deviceToken       string
		deviceType        string
		setupMocks        func(*mocks.MockQueries)
		shouldReturnError bool
		expectedResponse  *pb.DeviceToken
	}{
		{
			name:        "Success case",
			userID:      validUserID.String(),
			deviceToken: validDeviceToken,
			deviceType:  validDeviceType,
			setupMocks: func(db *mocks.MockQueries) {
				now := time.Now()
				returnedDeviceToken := database.DeviceToken{
					ID:          uuid.New(),
					UserID:      validUserID,
					DeviceToken: validDeviceToken,
					DeviceType:  validDeviceType,
					CreatedAt:   now,
					UpdatedAt:   now,
				}

				db.On("RegisterDeviceToken", mock.Anything, mock.MatchedBy(func(params database.RegisterDeviceTokenParams) bool {
					return params.UserID == validUserID &&
						params.DeviceToken == validDeviceToken &&
						params.DeviceType == validDeviceType
				})).Return(returnedDeviceToken, nil).Once()
			},
			shouldReturnError: false,
			expectedResponse: &pb.DeviceToken{
				UserId:      validUserID.String(),
				DeviceToken: validDeviceToken,
				DeviceType:  validDeviceType,
			},
		},
		{
			name:              "Invalid UUID",
			userID:            "not-a-valid-uuid",
			deviceToken:       validDeviceToken,
			deviceType:        validDeviceType,
			setupMocks:        func(*mocks.MockQueries) {},
			shouldReturnError: true,
			expectedResponse:  nil,
		},
		{
			name:        "Database error",
			userID:      validUserID.String(),
			deviceToken: validDeviceToken,
			deviceType:  validDeviceType,
			setupMocks: func(db *mocks.MockQueries) {
				db.On("RegisterDeviceToken", mock.Anything, mock.Anything).
					Return(database.DeviceToken{}, errors.New("database error")).Once()
			},
			shouldReturnError: true,
			expectedResponse:  nil,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh mocks for each test case
			mockDB := mocks.NewMockQueries()
			mockRabbitmq := mocks.NewMockRabbitMQ()
			mockFirebase := mocks.NewMockFirebaseClient()

			// Setup mocks according to test case
			tc.setupMocks(mockDB)

			// Create server with mocks
			srv := server.NewServer(mockDB, mockRabbitmq, "test/path", mockFirebase)

			// Make the request
			request := &pb.RegisterDeviceTokenRequest{
				UserId:      tc.userID,
				DeviceToken: tc.deviceToken,
				DeviceType:  tc.deviceType,
			}

			response, err := srv.RegisterDeviceToken(context.Background(), request)

			// Check error expectation
			if tc.shouldReturnError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotNil(t, response.DeviceToken)
				assert.Equal(t, tc.expectedResponse.UserId, response.DeviceToken.UserId)
				assert.Equal(t, tc.expectedResponse.DeviceToken, response.DeviceToken.DeviceToken)
				assert.Equal(t, tc.expectedResponse.DeviceType, response.DeviceToken.DeviceType)
				assert.NotNil(t, response.DeviceToken.CreatedAt)
				assert.NotNil(t, response.DeviceToken.UpdatedAt)
			}

			// Verify all expectations were met
			mockDB.AssertExpectations(t)
		})
	}
}
