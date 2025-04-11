[![CI](https://github.com/imhasandl/notification-service/actions/workflows/ci.yml/badge.svg)](https://github.com/imhasandl/notification-service/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/imhasandl/notification-service)](https://goreportcard.com/report/github.com/imhasandl/notification-service)
[![GoDoc](https://godoc.org/github.com/imhasandl/notification-service?status.svg)](https://godoc.org/github.com/imhasandl/notification-service)
[![Coverage](https://codecov.io/gh/imhasandl/notification-service/branch/main/graph/badge.svg)](https://codecov.io/gh/imhasandl/notification-service)
[![Go Version](https://img.shields.io/github/go-mod/go-version/imhasandl/notification-service)](https://golang.org/doc/devel/release.html)

# Notification Service

A microservice for handling notifications in a social media application, built with Go and gRPC.

## Overview

The Notification Service is responsible for sending and managing notifications for the social media platform. It uses gRPC for communication with other services and Firebase Cloud Messaging for delivering push notifications.

## Prerequisites

- Go 1.20 or later
- PostgreSQL database
- RabbitMQ
- Firebase project with Cloud Messaging enabled

## Configuration

Create a `.env` file in the root directory with the following variables:

```
PORT=":YOUR_GRPC_PORT"
DB_URL="postgres://username:password@host:port/database?sslmode=disable"
# DB_URL="postgres://username:password@db:port/database?sslmode=disable" // FOR DOCKER COMPOSE
TOKEN_SECRET="YOUR_JWT_SECRET_KEY"
RABBITMQ_URL="amqp://username:password@host:port"
FIREBASE_NOTIFICATION_KEY_PATH="path/to/firebase_key.json"
```

### Firebase Setup

1. Create a Firebase project at [console.firebase.google.com](https://console.firebase.google.com)
2. Navigate to Project Settings > Service Accounts
3. Generate a new private key (JSON format)
4. Save this file securely - it contains credentials to access your Firebase project
5. Update the `FIREBASE_NOTIFICATION_KEY_PATH` in your `.env` file to point to this JSON file

## Database Migrations

This service uses Goose for database migrations:

```bash
# Install Goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir migrations postgres "YOUR_DB_CONNECTION_STRING" up
```
## gRPC Methods

The service implements the following gRPC methods:

### SendNotification

Sends a push notification to a specific user.

#### Request Format

```json
{
   "title": "Notification title",
   "sender_username": "username of sender",
   "receiver_id": "UUID of recipient user",
   "content": "Notification message content", 
   "sent_at": "2023-01-01T12:00:00Z"
}
```

#### Response

```json
{
   "title": "Notification title",
   "sender_username": "Username of the sender",
   "receiver_id": "UUID of the recipient",
   "content": "The notification message",
   "sent_at": "Timestamp when notification was sent"
}
```

> **Note:** This method delivers notifications via Firebase Cloud Messaging if the user has a registered device token. If Firebase isn't initialized or no device token exists, the method will log this situation but still return a successful response.

### RegisterDeviceToken

Registers a user's device for push notifications.

#### Request Format

```json
{
   "user_id": "UUID of the user",
   "device_token": "Device-specific token for push notifications",
   "device_type": "Device platform (e.g., 'android', 'ios', 'web')"
}
```

#### Response Format

```json
{
   "device_token": {
      "id": "UUID of the device token record",
      "user_id": "UUID of the user",
      "device_token": "The device token string",
      "device_type": "Device platform type",
      "created_at": "Timestamp when the record was created",
      "updated_at": "Timestamp when the record was last updated"
   }
}
```

> **Note:** This method delivers notifications via Firebase Cloud Messaging if the user has a registered device token. If Firebase isn't initialized or no device token exists, the method will log this situation but still return a successful response.



### DeleteDeviceToken

Uses given ID of user's device token and deletes it in database.

#### Request Format

```json
{
   "user_id": "UUID of a user",
   "device_token": "user's device token"
}
```

#### Response Format

```json
{
   "status": "boolean value for the result of the request if status is TRUE request was successful and FALSE otherwise"
}
```

## RabbitMQ Integration

The Notification Service consumes messages from RabbitMQ to process asynchronous notification requests from other services.

### Message Consumption

The service automatically sets up and listens to:
- **Exchange**: `notifications.topic` (topic exchange)
- **Queue**: `notification_service_queue`
- **Routing Key**: `#` (wildcard - receives all messages published to the exchange)

### Publishing Messages to the Notification Service

Other microservices can send notification requests by publishing messages to the `notifications.topic` exchange. Messages should be JSON formatted with the following structure:

```json
{
   "title": "Notification title",
   "sender_id": "username of sender",
   "receiver_id": "UUID of recipient user",
   "content": "Notification message content", 
   "sent_at": "2023-01-01T12:00:00Z"
}
```

## Running the Service

```bash
go run cmd/main.go
```

## Docker Support

The service can be run as part of a Docker Compose setup along with other microservices. When using Docker, make sure to use the Docker Compose specific DB_URL configuration.
