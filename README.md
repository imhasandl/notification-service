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
## SendNotification

Sends a push notification to a specific user.

### Request Format

```json
{
   "title": "Notification title",
   "sender_username": "username of sender",
   "receiver_id": "UUID of recipient user",
   "content": "Notification message content", 
   "sent_at": "2023-01-01T12:00:00Z"
}
```

### Response

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
 * SendNotification sends a push notification to a specific user.
 *
 * Request format:
 * The method expects a notification payload in JSON format wrapped in a byte array:
 * {
 *   "title": "Notification title",
 *   "sender_username": "username of sender",
 *   "receiver_id": "UUID of recipient user",
 *   "content": "Notification message content",
 *   "sent_at": "2023-01-01T12:00:00Z" (timestamp in RFC3339 format)
 * }
 *
 * Response:
 * Returns a SendNotificationResponse containing:
 * - title: The notification title
 * - sender_username: Username of the sender
 * - receiver_id: UUID of the recipient
 * - content: The notification message
 * - sent_at: Timestamp when notification was sent
 *
 * Note: This method will attempt to deliver the notification via Firebase Cloud
 * Messaging if the user has a registered device token. If Firebase isn't
 * initialized or no device token is found, the method will log this situation
 * but will still return a successful response.
 */
## gRPC Methods

The service implements the following gRPC methods:

### RegisterDevice

Registers a user's device for push notifications.

### GetUserNotifications

Retrieves a list of notifications for a specific user.

### MarkNotificationAsRead

Marks a specific notification as read.

### DeleteNotification

Removes a notification from the user's list.

## Running the Service

```bash
go run cmd/main.go
```

## Docker Support

The service can be run as part of a Docker Compose setup along with other microservices. When using Docker, make sure to use the Docker Compose specific DB_URL configuration.