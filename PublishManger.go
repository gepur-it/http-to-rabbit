package main

import (
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type PublishManger struct {
	rabbit      *Rabbit
	isConnected bool
	mux         sync.Mutex
}

func (manager *PublishManger) Lock() {
	log.Print("request lock")
	manager.mux.Lock()
	log.Print("locked")
}

func (manager *PublishManger) Unlock() {
	log.Print("request unlock")
	manager.mux.Unlock()
	log.Print("unlocked")
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
	manager.Lock()
	defer manager.Unlock()

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
