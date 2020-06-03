package main

import (
	"log"

	"github.com/streadway/amqp"
)

type Rabbit struct {
	connectionString string
	connection       *amqp.Connection
	channel          *amqp.Channel
}

func (rabbit *Rabbit) Connect() error {
	conn, err := amqp.Dial(rabbit.connectionString)
	if err != nil {
		log.Printf("%s: %s", "Failed to connect to RabbitMQ", err)

		return err
	}

	rabbit.connection = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%s: %s", "Failed to create a channel", err)

		conn.Close()
		return err
	}

	rabbit.channel = ch

	return nil
}

func (rabbit *Rabbit) Close() {
	rabbit.channel.Close()
	rabbit.connection.Close()
}
