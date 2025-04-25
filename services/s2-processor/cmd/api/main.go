// cmd/api/main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/leandroalencar/banco-dados/shared/models"
	"github.com/leandroalencar/banco-dados/shared/utils"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	// Initialize RabbitMQ
	rabbitmq, err := utils.NewRabbitMQ(os.Getenv("MESSAGE_BROKER_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	// Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// User endpoints
	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)

	// Transaction endpoints
	r.POST("/transactions", createTransaction)
	r.GET("/transactions/:id", getTransaction)

	// Quotation endpoints
	r.GET("/quotations/latest", getLatestQuotation)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
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
