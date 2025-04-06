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

## RegisterDeviceToken

RegisterDeviceToken allows clients to register a device token for receiving push notifications. This method stores the association between a user and their device token in the database, enabling targeted delivery of notifications.

The method supports both new token registration and updating existing tokens. If a token for the user already exists, the record will be updated with the new information.

#### Request format

```json
{
   "user_id": "UUID of the user",
   "device_token": "Token provided by FCM or APNs",
   "device_type": "Device platform (e.g., android, ios, web)"
}
```

#### Response format

```json
{
   "device_token": {
      "id": "UUID of the registration entry",
      "user_id": "UUID of the user",
      "device_token": "The registered device token",
      "device_type": "Device platform",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
   }
}
```

> **Note:** Device tokens are unique per user and device. Attempting to register the same token for the same user will update the existing record rather than creating a duplicate.

## SendNotification

SendNotification processes notification requests that are consumed from RabbitMQ. Other microservices in the system can publish notification events to RabbitMQ queues, which are then picked up by the Notification Service. This allows for asynchronous notification handling across the platform.

The service listens to designated RabbitMQ queues and processes incoming messages according to their type and content. When a message is received, it is parsed and sent as a notification to the intended recipient.

#### Request format from different service

```json
{
   "title": "Notification title",
   "sender_username": "username of sender",
   "receiver_id": "UUID of recipient user",
   "content": "Notification message content", 
   "sent_at": "2023-01-01T12:00:00Z"
}
```

#### If message sent successfully the status response will be TRUE

```json
{
   "status": true
}
```

> **Note:** This method delivers notifications via Firebase Cloud Messaging if the user has a registered device token. If Firebase isn't initialized or no device token exists, the method will log this situation but still return a successful response.

## Running the Service

```bash
go run cmd/main.go
```

## Docker Support

The service can be run as part of a Docker Compose setup along with other microservices. When using Docker, make sure to use the Docker Compose specific DB_URL configuration.