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
		rabbit:      rabbit,
		isConnected: false,
	}

	return manager
}

func main() {
	log.Print("Start")

	config, err := configuration()
	if err != nil {
		log.Printf("%s: %s", "Configuration fail", err)
		return
	}

	manager := createPublishManger(config)

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

		err = manager.publish(queueName, string(b))
		if err != nil {
			log.Printf("%s", "Cant publish")
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}

		resp.Success()
	})

	log.Fatal(http.ListenAndServe(":7000", app))
}
