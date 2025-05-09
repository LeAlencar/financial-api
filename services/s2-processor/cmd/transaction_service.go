package main

import (
	"encoding/json"
	"log"

	"github.com/gocql/gocql"
	"github.com/leandroalencar/banco-dados/shared/models"
	"github.com/streadway/amqp"
)

type TransactionService struct {
	channel *amqp.Channel
	session *gocql.Session
}

func NewTransactionService(channel *amqp.Channel) *TransactionService {
	// Connect to Cassandra
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "transactions"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}

	// Create keyspace and table if they don't exist
	if err := session.Query(`
		CREATE KEYSPACE IF NOT EXISTS transactions
		WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}
	`).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}

	if err := session.Query(`
		CREATE TABLE IF NOT EXISTS transactions.transactions (
			id text PRIMARY KEY,
			user_id text,
			type text,
			currency_pair text,
			amount double,
			exchange_rate double,
			total_value double,
			status text,
			timestamp timestamp,
			quotation_id text
		)
	`).Exec(); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return &TransactionService{
		channel: channel,
		session: session,
	}
}

func (s *TransactionService) Start() {
	// Declare queue
	_, err := s.channel.QueueDeclare(
		"transactions", // queue name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Consume messages
	msgs, err := s.channel.Consume(
		"transactions", // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var transaction models.Transaction
			if err := json.Unmarshal(msg.Body, &transaction); err != nil {
				log.Printf("Error unmarshaling transaction: %v", err)
				continue
			}

			// Insert transaction into Cassandra
			err := s.session.Query(`
				INSERT INTO transactions.transactions (
					id, user_id, type, currency_pair, amount, exchange_rate,
					total_value, status, timestamp, quotation_id
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, transaction.ID, transaction.UserID, transaction.Type,
				transaction.CurrencyPair, transaction.Amount, transaction.ExchangeRate,
				transaction.TotalValue, transaction.Status, transaction.Timestamp,
				transaction.QuotationID).Exec()

			if err != nil {
				log.Printf("Error saving transaction: %v", err)
			} else {
				log.Printf("Transaction processed: %s", transaction.ID)
			}
		}
	}()

	<-forever
}
