package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/domain/repositories"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/domain/services"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/infra/database"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/infra/messaging"
	"github.com/leandroalencar/banco-dados/shared/utils"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize PostgreSQL connection
	postgresDB, err := database.NewPostgresConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer postgresDB.Close()

	// Initialize MongoDB connection
	mongoDB, err := database.NewMongoConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository(postgresDB.GetPool())
	userService := services.NewUserService(userRepo)

	quotationRepo := repositories.NewQuotationRepository(mongoDB.GetDatabase())
	quotationService := services.NewQuotationService(quotationRepo)

	// Initialize currency transaction repository
	currencyTransactionRepo := repositories.NewCurrencyTransactionRepository(mongoDB.GetDatabase())

	// Initialize RabbitMQ connection URI
	rabbitURI := os.Getenv("RABBITMQ_URI")
	if rabbitURI == "" {
		rabbitURI = "amqp://guest:guest@localhost:5672/"
	}

	// Initialize RabbitMQ for transaction service
	rabbitmq, err := utils.NewRabbitMQ(rabbitURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	// Initialize transaction service
	transactionService := services.NewTransactionService(userRepo, quotationRepo, currencyTransactionRepo, rabbitmq)

	// Create consumers
	userConsumer := messaging.NewUserConsumer(rabbitURI, userService)
	quotationConsumer := messaging.NewQuotationConsumer(rabbitURI, quotationService)
	transactionConsumer := messaging.NewTransactionConsumer(rabbitURI, transactionService)

	// Create consumer manager
	manager := messaging.NewConsumerManager(userConsumer, quotationConsumer, transactionConsumer)

	// Channel to handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start all consumers
	if err := manager.StartAll(ctx); err != nil {
		log.Fatalf("Failed to start consumers: %v", err)
	}

	log.Println("All consumers started. Press CTRL+C to stop.")

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down gracefully...")

	// Cancel context to stop consumers
	cancel()

	// Wait for all consumers to finish
	manager.WaitForShutdown()
	log.Println("All consumers stopped. Goodbye!")
}
