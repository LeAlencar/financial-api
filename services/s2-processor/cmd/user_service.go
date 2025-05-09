package main

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/leandroalencar/banco-dados/shared/models"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

type UserService struct {
	channel *amqp.Channel
	db      *sql.DB
}

func NewUserService(channel *amqp.Channel) *UserService {
	// Connect to PostgreSQL
	db, err := sql.Open("postgres", "postgres://admin:password@localhost:5432/users?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Create users table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			balance DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	return &UserService{
		channel: channel,
		db:      db,
	}
}

func (s *UserService) Start() {
	// Declare queue
	_, err := s.channel.QueueDeclare(
		"users", // queue name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Consume messages
	msgs, err := s.channel.Consume(
		"users", // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var user models.User
			if err := json.Unmarshal(msg.Body, &user); err != nil {
				log.Printf("Error unmarshaling user: %v", err)
				continue
			}

			// Insert or update user in PostgreSQL
			_, err := s.db.Exec(`
				INSERT INTO users (id, name, email, balance, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (id) DO UPDATE
				SET name = $2, email = $3, balance = $4, updated_at = $6
			`, user.ID, user.Name, user.Email, user.Balance, user.CreatedAt, user.UpdatedAt)

			if err != nil {
				log.Printf("Error saving user: %v", err)
			} else {
				log.Printf("User processed: %s", user.ID)
			}
		}
	}()

	<-forever
}
