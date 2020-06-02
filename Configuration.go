package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	ConnectionString string `json:"rabbit_connection_string"`
}

func configuration() (Configuration, error) {
	file, Err := os.Open("config.json")
	if Err != nil {
		fmt.Println("error while loading config from config.json :", Err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	return configuration, err
}
