package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Content   string    `bson:"content"`
	Timestamp time.Time `bson:"timestamp"`
}

func main() {
	// Conectar ao MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Erro ao conectar no MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Verificar conexão com MongoDB
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Erro ao verificar conexão com MongoDB: %v", err)
	}

	// Obter coleção
	collection := client.Database("financial").Collection("messages")

	// Conectar ao RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Erro ao conectar no RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir canal: %v", err)
	}
	defer ch.Close()

	queueName := "cotacoes"

	// Declarar a fila (precaução, caso ela ainda não exista)
	_, err = ch.QueueDeclare(
		queueName, // nome
		false,     // durável
		false,     // autoDelete
		false,     // exclusiva
		false,     // noWait
		nil,       // argumentos
	)
	if err != nil {
		log.Fatalf("Erro ao declarar fila: %v", err)
	}

	msgs, err := ch.Consume(
		queueName, // nome da fila
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Erro ao registrar consumidor: %v", err)
	}

	fmt.Println("Aguardando mensagens. Para sair pressione CTRL+C")

	// Loop para receber mensagens
	for msg := range msgs {
		message := Message{
			Content:   string(msg.Body),
			Timestamp: time.Now(),
		}

		// Salvar no MongoDB
		_, err := collection.InsertOne(ctx, message)
		if err != nil {
			log.Printf("Erro ao salvar mensagem no MongoDB: %v", err)
			continue
		}

		fmt.Printf("Mensagem recebida e salva: %s\n", msg.Body)
	}
}
