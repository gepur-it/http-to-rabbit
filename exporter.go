package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

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
	log.Print("Start")

	config, err := configuration()
	if err != nil {
		log.Printf("%s: %s", "Configuration fail", err)
		return
	}

	manager := createPublishManger(config)

	messageChanel := make(chan Message)
	defer close(messageChanel)

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

	app := &App{
		DefaultRoute: func(resp Response, req Request) {
			resp.Text(http.StatusNotFound, "Not found")
		},
	}

	app.Handle(`^/([\w\._-]+)$`, func(resp Response, req Request) {
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
			log.Printf("%s", "Cant publish")
			http.Error(resp, PublishResult.error, http.StatusInternalServerError)
			return
		}

		resp.Success()
	})

	log.Fatal(http.ListenAndServe(":7000", app))
}
