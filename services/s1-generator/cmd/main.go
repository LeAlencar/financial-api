package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/leandroalencar/banco-dados/shared/models"
	"github.com/leandroalencar/banco-dados/shared/utils"
)

func main() {
	// Initialize RabbitMQ connection
	rabbitmq, err := utils.NewRabbitMQ("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	// Start generating messages
	for {
		// Generate mock quotation
		quotation := generateMockQuotation()
		err = rabbitmq.PublishMessage("quotations", quotation)
		if err != nil {
			log.Printf("Error publishing quotation: %v", err)
		}

		// Generate mock transaction
		transaction := generateMockTransaction()
		err = rabbitmq.PublishMessage("transactions", transaction)
		if err != nil {
			log.Printf("Error publishing transaction: %v", err)
		}

		time.Sleep(5 * time.Second)
	}
}

func generateMockQuotation() models.Quotation {
	return models.Quotation{
		CurrencyPair:  "USD/BRL",
		BuyPrice:      4.90 + rand.Float64()*0.10,
		SellPrice:     5.00 + rand.Float64()*0.10,
		Timestamp:     time.Now(),
		LastUpdatedBy: "generator",
	}
}

func generateMockTransaction() models.Transaction {
	transactionTypes := []models.TransactionType{models.Buy, models.Sell}
	return models.Transaction{
		UserID:       "user-" + string(rand.Intn(100)),
		Type:         transactionTypes[rand.Intn(2)],
		CurrencyPair: "USD/BRL",
		Amount:       float64(rand.Intn(1000)) + rand.Float64(),
		Status:       "PENDING",
		Timestamp:    time.Now(),
	}
}
