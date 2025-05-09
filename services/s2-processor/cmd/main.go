// cmd/api/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/leandroalencar/banco-dados/shared/models"
	"github.com/streadway/amqp"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create channels for each service
	userChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open user channel: %v", err)
	}
	defer userChan.Close()

	quotationChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open quotation channel: %v", err)
	}
	defer quotationChan.Close()

	transactionChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open transaction channel: %v", err)
	}
	defer transactionChan.Close()

	// Initialize services
	userService := NewUserService(userChan)
	quotationService := NewQuotationService(quotationChan)
	transactionService := NewTransactionService(transactionChan)

	// Start services
	go userService.Start()
	go quotationService.Start()
	go transactionService.Start()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down services...")
}

func createUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement user creation logic

	c.JSON(http.StatusCreated, user)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement user retrieval logic

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func createTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement transaction creation logic

	c.JSON(http.StatusCreated, transaction)
}

func getTransaction(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement transaction retrieval logic

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func getLatestQuotation(c *gin.Context) {
	// TODO: Implement latest quotation retrieval logic

	c.JSON(http.StatusOK, gin.H{
		"currency_pair": "USD/BRL",
		"buy_price":     4.95,
		"sell_price":    5.05,
	})
}
