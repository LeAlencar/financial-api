package main

import (
	"encoding/json"
	"log"

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

	// Handle quotation messages
	err = rabbitmq.ConsumeMessages("quotations", func(msg []byte) error {
		var quotation models.Quotation
		if err := json.Unmarshal(msg, &quotation); err != nil {
			return err
		}
		log.Printf("Validated quotation: %+v\n", quotation)
		return nil
	})
	if err != nil {
		log.Printf("Error consuming quotations: %v", err)
	}

	// Handle transaction messages
	err = rabbitmq.ConsumeMessages("transactions", func(msg []byte) error {
		var transaction models.Transaction
		if err := json.Unmarshal(msg, &transaction); err != nil {
			return err
		}
		log.Printf("Validated transaction: %+v\n", transaction)
		return nil
	})
	if err != nil {
		log.Printf("Error consuming transactions: %v", err)
	}

	// Keep the service running
	select {}
}
