// cmd/api/main.go
package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/domain/repositories"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/domain/services"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/infra/database"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/infra/http/routes"
	"github.com/leandroalencar/banco-dados/shared/utils"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	ctx := context.Background()

	// Initialize Database
	db, err := database.NewPostgresConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository(db.GetPool())
	userService := services.NewUserService(userRepo)

	// Initialize RabbitMQ
	rabbitmq, err := utils.NewRabbitMQ(os.Getenv("MESSAGE_BROKER_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	// Setup router
	router := routes.SetupRouter(userService)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
