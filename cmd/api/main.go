// cmd/api/main.go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lealencar/financial-api/internal/domain/repositories"
	"github.com/lealencar/financial-api/internal/infra/api/awesomeapi"
	"github.com/lealencar/financial-api/internal/infra/database"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	// Conexão com PostgreSQL (existente)
	_, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Nova conexão com MongoDB
	mongoDB, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoDB.Client().Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Inicializar dependências
	currencyRepo := repositories.NewCurrencyRepository(mongoDB)
	apiClient := awesomeapi.NewClient()

	r := gin.Default()

	// Nova rota para buscar e salvar cotações
	r.GET("/api/exchange", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// Buscar cotação USD-BRL
		currency, err := apiClient.GetExchangeRate("USD", "BRL")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch exchange rate", "details": err.Error()})
			return
		}

		// Salvar no MongoDB
		if err := currencyRepo.Insert(ctx, currency); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save data", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message":  "Exchange rate saved successfully",
			"currency": currency,
		})
	})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
