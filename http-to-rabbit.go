package main

import (
	"io/ioutil"
	"log"
	"net/http"

	logger "github.com/sirupsen/logrus"

	application "gepur.com/http-to-rabbit/application"
)

type Message struct {
	queueName   string
	body        string
	errorChanel chan<- PublishResult
}

type PublishResult struct {
	success bool
	error   string
}

func main() {
	config, err := configuration()
	if err != nil {
		log.Printf("%s: %s", "Configuration fail", err)
		return
	}

	initLogger(config)

	manager := createPublishManger(config)

	messageChanel := make(chan Message)
	defer close(messageChanel)

	logger.Info("Service strated")

	// register publisher worker
	go func(messages <-chan Message) {
		for {
			msg := <-messages

			err = manager.publishWithReconnects(msg.queueName, msg.body)

			var res PublishResult
			if err != nil {
				res = PublishResult{
					success: false,
					error:   err.Error(),
				}
			} else {
				res = PublishResult{
					success: true,
					error:   "",
				}
			}

			msg.errorChanel <- res
		}
	}(messageChanel)

	app := &application.App{
		DefaultRoute: func(resp application.Response, req application.Request) {
			resp.Text(http.StatusNotFound, "Not found")
		},
	}

	app.Handle(`^/([\w\._-]+)$`, func(resp application.Response, req application.Request) {
		queueName := req.Params[0]

		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()

		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}

		respChan := make(chan PublishResult)
		defer close(respChan)

		messageChanel <- Message{
			queueName:   queueName,
			body:        string(b),
			errorChanel: respChan,
		}

		PublishResult := <-respChan
		if !PublishResult.success {
			resp.Error(PublishResult.error)
			return
		}

		resp.Success()
	})

	logger.Fatal(http.ListenAndServe(":7000", app))
}
