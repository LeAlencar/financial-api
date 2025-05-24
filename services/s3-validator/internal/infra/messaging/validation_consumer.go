package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/leandroalencar/banco-dados/services/s3-validator/internal/domain/services"
	"github.com/leandroalencar/banco-dados/shared/messaging/events"
	"github.com/streadway/amqp"
)

type ValidationConsumer struct {
	uri               string
	validationService *services.ValidationService
	queueName         string
	eventType         string
}

func NewUserValidationConsumer(uri string, validationService *services.ValidationService) *ValidationConsumer {
	return &ValidationConsumer{
		uri:               uri,
		validationService: validationService,
		queueName:         "users-validator",
		eventType:         "user",
	}
}

func NewTransactionValidationConsumer(uri string, validationService *services.ValidationService) *ValidationConsumer {
	return &ValidationConsumer{
		uri:               uri,
		validationService: validationService,
		queueName:         "transactions-validator",
		eventType:         "transaction",
	}
}

func NewQuotationValidationConsumer(uri string, validationService *services.ValidationService) *ValidationConsumer {
	return &ValidationConsumer{
		uri:               uri,
		validationService: validationService,
		queueName:         "quotations-validator",
		eventType:         "quotation",
	}
}

func (c *ValidationConsumer) Name() string {
	return fmt.Sprintf("%s_validation_consumer", c.eventType)
}

func (c *ValidationConsumer) Start(ctx context.Context) error {
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

	_, err = ch.QueueDeclare(
		c.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error declaring queue: %v", err)
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("error setting QoS: %v", err)
	}

	msgs, err := ch.Consume(
		c.queueName,
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

	log.Printf("%s validation consumer started. Listening for messages on queue: %s", c.eventType, c.queueName)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("message channel closed")
			}

			if err := c.processMessage(ctx, msg); err != nil {
				log.Printf("Error processing %s validation message: %v", c.eventType, err)
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)
		}
	}
}

func (c *ValidationConsumer) processMessage(ctx context.Context, msg amqp.Delivery) error {
	rawPayload := string(msg.Body)

	switch c.eventType {
	case "user":
		var event events.UserEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			return fmt.Errorf("error decoding user event: %v", err)
		}

		if err := c.validationService.ValidateUserEvent(ctx, &event, rawPayload); err != nil {
			return fmt.Errorf("error validating user event: %v", err)
		}

		log.Printf("Successfully validated user event: %s for user ID: %d", event.Action, event.Data.ID)

	case "transaction":
		var event events.TransactionEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			return fmt.Errorf("error decoding transaction event: %v", err)
		}

		if err := c.validationService.ValidateTransactionEvent(ctx, &event, rawPayload); err != nil {
			return fmt.Errorf("error validating transaction event: %v", err)
		}

		log.Printf("Successfully validated transaction event: %s for user ID: %s", event.Action, event.Data.UserID)

	case "quotation":
		var event events.QuotationEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			return fmt.Errorf("error decoding quotation event: %v", err)
		}

		if err := c.validationService.ValidateQuotationEvent(ctx, &event, rawPayload); err != nil {
			return fmt.Errorf("error validating quotation event: %v", err)
		}

		log.Printf("Successfully validated quotation event: %s for currency pair: %s", event.Action, event.Data.CurrencyPair)

	default:
		return fmt.Errorf("unknown event type: %s", c.eventType)
	}

	return nil
}
