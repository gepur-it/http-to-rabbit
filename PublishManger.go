package main

import (
	"fmt"

	logger "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type PublishManger struct {
	rabbit       *Rabbit
	retriesCount int
	isConnected  bool
}

func createPublishManger(config Configuration) *PublishManger {
	rabbit := &Rabbit{
		config.ConnectionString,
		nil,
		nil,
	}

	manager := &PublishManger{
		rabbit:       rabbit,
		retriesCount: config.RetriesCount,
		isConnected:  false,
	}

	return manager
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
		return err
	}

	err = manager.rabbit.channel.Publish(
		queue, // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		manager.Disconnect()
		return err
	}

	err = manager.rabbit.channel.TxCommit()
	if err != nil {
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
		} else {
			logger.WithFields(logger.Fields{
				"error": err,
				"queue": queue,
				"body":  body,
			}).Error(fmt.Sprintf("Failed to publish a message. Аttempt #%d", i+1))
		}
	}

	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err,
			"queue": queue,
			"body":  body,
		}).Error("Failed to publish a message. Maximum attempts reached")
	}

	return err
}
