package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/domain/repositories"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/domain/services"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/infra/database"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/infra/http/routes"
	"github.com/leandroalencar/banco-dados/shared/utils"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Initialize Database
	ctx := context.Background()
	db, err := database.NewPostgresConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository(db.GetPool())
	userService := services.NewUserService(userRepo)

	// Initialize RabbitMQ
	rabbitURI := os.Getenv("RABBITMQ_URI")
	if rabbitURI == "" {
		rabbitURI = "amqp://guest:guest@localhost:5672/"
	}

	rabbitmq, err := utils.NewRabbitMQ(rabbitURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	// Start HTTP server
	router := routes.SetupRouter(rabbitmq, userService)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("HTTP server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
