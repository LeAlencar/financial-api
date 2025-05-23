package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/domain/services"
	"github.com/leandroalencar/banco-dados/shared/messaging/events"
	"github.com/streadway/amqp"
)

type TransactionConsumer struct {
	uri                string
	transactionService *services.TransactionService
}

func NewTransactionConsumer(uri string, transactionService *services.TransactionService) *TransactionConsumer {
	return &TransactionConsumer{
		uri:                uri,
		transactionService: transactionService,
	}
}

// Name returns the consumer's name
func (c *TransactionConsumer) Name() string {
	return "transaction_consumer"
}

// Start implements the Consumer interface
func (c *TransactionConsumer) Start(ctx context.Context) error {
	// Connect to RabbitMQ
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

	queueName := "transactions"

	// Declare the queue
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("error declaring queue: %v", err)
	}

	// Set QoS
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("error setting QoS: %v", err)
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
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

			var event events.TransactionEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Error decoding message: %v", err)
				msg.Nack(false, true)
				continue
			}

			if err := c.transactionService.ProcessTransactionEvent(ctx, &event); err != nil {
				log.Printf("Error processing transaction event: %v", err)
				msg.Nack(false, true)
				continue
			}

			log.Printf("Successfully processed transaction event: %s for user ID: %s", event.Action, event.Data.UserID)
			msg.Ack(false)
		}
	}
}
