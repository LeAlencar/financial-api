package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocql/gocql"
	"github.com/leandroalencar/banco-dados/services/s3-validator/internal/domain/repositories"
	"github.com/leandroalencar/banco-dados/services/s3-validator/internal/domain/services"
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
	validationLogRepo := repositories.NewValidationLogRepository(session)

	// Initialize service
	validationService := services.NewValidationService(validationLogRepo, transactionRepo)

	// Initialize RabbitMQ connection
	rabbitmqURI := "amqp://guest:guest@localhost:5672/"

	// Initialize validation consumers
	userValidationConsumer := messaging.NewUserValidationConsumer(rabbitmqURI, validationService)
	transactionValidationConsumer := messaging.NewTransactionValidationConsumer(rabbitmqURI, validationService)
	quotationValidationConsumer := messaging.NewQuotationValidationConsumer(rabbitmqURI, validationService)

	// Create consumer manager
	manager := messaging.NewConsumerManager(
		userValidationConsumer,
		transactionValidationConsumer,
		quotationValidationConsumer,
	)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start all consumers
	if err := manager.StartAll(ctx); err != nil {
		log.Fatalf("Failed to start consumers: %v", err)
	}

	log.Println("All validation consumers started. Press CTRL+C to stop.")

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down gracefully...")

	// Cancel context to stop consumers
	cancel()

	// Wait for all consumers to finish
	manager.WaitForShutdown()
	log.Println("All consumers stopped. Goodbye!")
}
