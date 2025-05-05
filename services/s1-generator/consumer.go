package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
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
		fmt.Printf("Mensagem recebida: %s\n", msg.Body)
	}
}
