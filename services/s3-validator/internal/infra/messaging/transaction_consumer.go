package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/leandroalencar/banco-dados/services/s3-validator/internal/domain/models"
	"github.com/leandroalencar/banco-dados/services/s3-validator/internal/domain/repositories"
	"github.com/streadway/amqp"
)

type TransactionConsumer struct {
	uri             string
	transactionRepo *repositories.TransactionRepository
}

func NewTransactionConsumer(uri string, transactionRepo *repositories.TransactionRepository) *TransactionConsumer {
	return &TransactionConsumer{
		uri:             uri,
		transactionRepo: transactionRepo,
	}
}

func (c *TransactionConsumer) Name() string {
	return "transaction_consumer"
}

func (c *TransactionConsumer) Start(ctx context.Context) error {
	conn, err := amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("error connecting to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("error opening channel: %v", err)
	}
	defer ch.Close()

	queueName := "transactions-validator"

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error declaring queue: %v", err)
	}

	err = ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		return fmt.Errorf("error setting QoS: %v", err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error registering consumer: %v", err)
	}

	log.Printf("Transaction consumer started. Listening for messages on queue: %s", queueName)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("message channel closed")
			}

			var transaction models.Transaction
			if err := json.Unmarshal(msg.Body, &transaction); err != nil {
				log.Printf("Error decoding message: %v", err)
				msg.Nack(false, true)
				continue
			}

			if err := c.transactionRepo.Save(ctx, &transaction); err != nil {
				log.Printf("Error saving transaction: %v", err)
				msg.Nack(false, true)
				continue
			}

			log.Printf("Successfully saved transaction: %s for user: %s with status: %s",
				transaction.ID, transaction.UserID, transaction.Status)
			msg.Ack(false)
		}
	}
}
