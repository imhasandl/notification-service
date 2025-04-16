package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/imhasandl/notification-service/cmd/server"
	"github.com/imhasandl/notification-service/internal/database"
	"github.com/imhasandl/notification-service/internal/firebase"
	"github.com/imhasandl/notification-service/internal/rabbitmq"
	pb "github.com/imhasandl/notification-service/protos"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the postgres driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Config holds application configuration
type Config struct {
	port            string
	dbURL           string
	rabbitmqURL     string
	firebaseKeyPath string
}

func main() {
	// Load config
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup listeners and connections
	listener, err := net.Listen("tcp", config.port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Initialize services
	dbQueries, dbConn, err := initDatabase(config.dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbConn.Close()

	rmq, err := initRabbitMQ(config.rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rmq.Close()

	fb, err := initFirebase(context.Background(), config.firebaseKeyPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Create and start server
	srv := server.NewServer(dbQueries, rmq, config.firebaseKeyPath, fb)
	startServer(listener, srv)
}

// loadConfig loads configuration from environment variables
func loadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("PORT environment variable not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL environment variable not set")
	}

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL environment variable not set")
	}

	firebaseKeyPath := os.Getenv("FIREBASE_NOTIFICATION_KEY_PATH")
	if firebaseKeyPath == "" {
		return nil, fmt.Errorf("FIREBASE_NOTIFICATION_KEY_PATH environment variable not set")
	}

	return &Config{
		port:            port,
		dbURL:           dbURL,
		rabbitmqURL:     rabbitmqURL,
		firebaseKeyPath: firebaseKeyPath,
	}, nil
}

// initDatabase initializes the database connection
func initDatabase(dbURL string) (*database.Queries, *sql.DB, error) {
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, err
	}
	dbQueries := database.New(dbConn)
	return dbQueries, dbConn, nil
}

// initRabbitMQ initializes the RabbitMQ client
func initRabbitMQ(rabbitmqURL string) (*rabbitmq.RabbitMQ, error) {
	return rabbitmq.NewRabbitMQ(rabbitmqURL)
}

// initFirebase initializes the Firebase client
func initFirebase(ctx context.Context, keyPath string) (*firebase.Client, error) {
	return firebase.InitFirebase(ctx, keyPath)
}

// startServer initializes and starts the gRPC server
func startServer(lis net.Listener, srv *server.Server) {
	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	log.Printf("Server listening on %v", lis.Addr())

	// Start consuming messages from notification-queue
	go srv.Consume()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
