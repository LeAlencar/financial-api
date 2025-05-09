package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/leandroalencar/banco-dados/shared/models"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QuotationService struct {
	channel *amqp.Channel
	client  *mongo.Client
	db      *mongo.Database
}

func NewQuotationService(channel *amqp.Channel) *QuotationService {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database("quotations")

	return &QuotationService{
		channel: channel,
		client:  client,
		db:      db,
	}
}

func (s *QuotationService) Start() {
	// Declare queue
	_, err := s.channel.QueueDeclare(
		"quotations", // queue name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Consume messages
	msgs, err := s.channel.Consume(
		"quotations", // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var quotation models.Quotation
			if err := json.Unmarshal(msg.Body, &quotation); err != nil {
				log.Printf("Error unmarshaling quotation: %v", err)
				continue
			}

			// Insert quotation into MongoDB
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := s.db.Collection("quotations").InsertOne(ctx, quotation)
			cancel()

			if err != nil {
				log.Printf("Error saving quotation: %v", err)
			} else {
				log.Printf("Quotation processed: %s", quotation.ID)
			}
		}
	}()

	<-forever
}
