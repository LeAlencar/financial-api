package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocql/gocql"
	"github.com/leandroalencar/banco-dados/services/s3-validator/internal/domain/repositories"
	"github.com/leandroalencar/banco-dados/services/s3-validator/internal/infra/messaging"
)

func main() {
	// Initialize Cassandra connection
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "financial"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	// Initialize repositories
	transactionRepo := repositories.NewTransactionRepository(session)

	// Initialize RabbitMQ connection
	rabbitmqURI := "amqp://guest:guest@localhost:5672/"

	// Initialize consumers
	transactionConsumer := messaging.NewTransactionConsumer(rabbitmqURI, transactionRepo)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start consumers
	if err := transactionConsumer.Start(ctx); err != nil {
		log.Printf("Error starting transaction consumer: %v", err)
	}

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down...")
	cancel()
}
