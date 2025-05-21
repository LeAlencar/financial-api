package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/domain/services"
	"github.com/leandroalencar/banco-dados/shared/models"
	"github.com/streadway/amqp"
)

type QuotationConsumer struct {
	uri              string
	quotationService *services.QuotationService
}

func NewQuotationConsumer(uri string, quotationService *services.QuotationService) *QuotationConsumer {
	return &QuotationConsumer{
		uri:              uri,
		quotationService: quotationService,
	}
}

// Name returns the consumer's name
func (c *QuotationConsumer) Name() string {
	return "quotation_consumer"
}

// Start implements the Consumer interface
func (c *QuotationConsumer) Start(ctx context.Context) error {
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

	queueName := "quotations"

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

	log.Printf("Quotation consumer started. Listening for messages on queue: %s", queueName)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("message channel closed")
			}

			var quotation models.Quotation
			if err := json.Unmarshal(msg.Body, &quotation); err != nil {
				log.Printf("Error decoding message: %v", err)
				msg.Nack(false, true)
				continue
			}

			if err := c.quotationService.SaveQuotation(ctx, &quotation); err != nil {
				log.Printf("Error saving quotation: %v", err)
				msg.Nack(false, true)
				continue
			}

			log.Printf("Successfully saved quotation for pair: %s at price: %.2f/%.2f",
				quotation.CurrencyPair, quotation.BuyPrice, quotation.SellPrice)
			msg.Ack(false)
		}
	}
}
