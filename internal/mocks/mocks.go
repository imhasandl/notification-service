package mocks

import (
	"context"

	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"github.com/imhasandl/notification-service/internal/database"
	"github.com/imhasandl/notification-service/internal/firebase"
	"github.com/imhasandl/notification-service/internal/rabbitmq"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

// MockQueries mocks the database.Queries struct
type MockQueries struct {
	mock.Mock
}

// NewMockQueries creates and returns a new mock database queries object
func NewMockQueries() *MockQueries {
	return &MockQueries{}
}

// RegisterDeviceToken mocks the database method for registering device tokens
func (m *MockQueries) RegisterDeviceToken(ctx context.Context, arg database.RegisterDeviceTokenParams) (database.DeviceToken, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.DeviceToken), args.Error(1)
}

// GetDeviceTokensByUserID mocks the database method for fetching device tokens
func (m *MockQueries) GetDeviceTokensByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

// DeleteDeviceToken mocks the database method for deleting device tokens
func (m *MockQueries) DeleteDeviceToken(ctx context.Context, arg database.DeleteDeviceTokenParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// SendNotification mocks the database method for sending notifications
func (m *MockQueries) SendNotification(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockFirebaseClient mocks the Firebase client for sending notifications
type MockFirebaseClient struct {
	mock.Mock
	FCMClient *MockFCMClient // Add this field to match the real FirebaseClient structure
}

// NewMockFirebaseClient creates a new mock firebase client
func NewMockFirebaseClient() *MockFirebaseClient {
	fcmClient := &MockFCMClient{}
	return &MockFirebaseClient{
		FCMClient: fcmClient,
	}
}

// GetMessagingClient returns the FCM client for sending messages
func (m *MockFirebaseClient) GetMessagingClient() firebase.MessagingClient {
	args := m.Called()
	if mockClient := args.Get(0); mockClient != nil {
		return mockClient.(firebase.MessagingClient)
	}
	return m.FCMClient
}

// MockFCMClient mocks the Firebase Cloud Messaging client
type MockFCMClient struct {
	mock.Mock
}

// Send mocks the FCM Send method
func (m *MockFCMClient) Send(ctx context.Context, message *messaging.Message) (string, error) {
	// Add nil check to prevent panic
	if m == nil {
		return "", nil
	}

	args := m.Called(ctx, message)
	return args.String(0), args.Error(1)
}

// Update MockFCMClient to implement MessagingClient
var _ firebase.MessagingClient = (*MockFCMClient)(nil)

// MockRabbitMQ mocks the RabbitMQ client
type MockRabbitMQ struct {
	mock.Mock
	Channel *MockChannel
}

// NewMockRabbitMQ creates a new mock RabbitMQ client
func NewMockRabbitMQ() *MockRabbitMQ {
	channel := &MockChannel{}
	return &MockRabbitMQ{
		Channel: channel,
	}
}

// Close mocks the RabbitMQ Close method
func (m *MockRabbitMQ) Close() {
	m.Called()
}

// Ensure MockRabbitMQ implements the RabbitMQClient interface
var _ rabbitmq.Client = (*MockRabbitMQ)(nil)

// GetChannel returns the mock channel
func (m *MockRabbitMQ) GetChannel() *amqp.Channel {
	args := m.Called()
	return args.Get(0).(*amqp.Channel)
}

// MockChannel mocks the RabbitMQ Channel
type MockChannel struct {
	mock.Mock
}

// Consume mocks the RabbitMQ Consume method
func (m *MockChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args map[string]interface{}) (<-chan messaging.Message, error) {
	mockArgs := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	return mockArgs.Get(0).(<-chan messaging.Message), mockArgs.Error(1)
}

// MockDBQuerier is a mock for the DBQuerier interface
type MockDBQuerier struct {
	mock.Mock
}

func (m *MockDBQuerier) GetDeviceTokensByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockDBQuerier) RegisterDeviceToken(ctx context.Context, arg database.RegisterDeviceTokenParams) (database.DeviceToken, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.DeviceToken), args.Error(1)
}

func (m *MockDBQuerier) DeleteDeviceToken(ctx context.Context, arg database.DeleteDeviceTokenParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockDBQuerier) SendNotification(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockRabbitMQClient is a mock for the RabbitMQ Client interface
type MockRabbitMQClient struct {
	mock.Mock
}

func (m *MockRabbitMQClient) Close() {
	m.Called()
}

func (m *MockRabbitMQClient) GetChannel() *amqp.Channel {
	args := m.Called()
	return args.Get(0).(*amqp.Channel)
}
