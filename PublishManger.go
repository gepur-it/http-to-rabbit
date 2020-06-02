package main

import (
	"log"

	"github.com/streadway/amqp"
)

type PublishManger struct {
	rabbit       *Rabbit
	retriesCount int
	isConnected  bool
}

func (manager *PublishManger) Connect() error {
	if manager.isConnected == true {
		return nil
	}

	err := manager.rabbit.Connect()
	if err != nil {
		return err
	}

	manager.isConnected = true

	return nil
}

func (manager *PublishManger) Disconnect() {
	if manager.isConnected == false {
		return
	}

	manager.rabbit.Close()
	manager.isConnected = false
}

func (manager *PublishManger) publish(queue string, body string) error {
	err := manager.Connect()
	if err != nil {
		log.Printf("%s", "Connection error")

		return err
	}

	q, err := manager.rabbit.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Printf("%s", "Error declare")
		manager.Disconnect()

		return err
	}

	err = manager.rabbit.channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		log.Printf("Failed to publish a message: %s", err)
		manager.Disconnect()
	}

	return err
}

func (manager *PublishManger) publishWithReconnects(queue string, body string) error {
	var err error
	for i := 0; i <= manager.retriesCount; i++ {
		err = manager.publish(queue, body)
		if err == nil {
			break
		}
	}

	return err
}
