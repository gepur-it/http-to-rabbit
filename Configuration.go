package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type LogstashConfig struct {
	ApplicationName string `json:"application_name"`
	Protocol        string `json:"protocol"`
	Address         string `json:"address"`
}

type EmailConfig struct {
	ApplicationName string `json:"application_name"`
	SmtpHost        string `json:"smtp_host"`
	SmtpPort        int    `json:"smtp_port"`
	SmtpFrom        string `json:"smtp_from"`
	SmtpTo          string `json:"smtp_to"`
	SmtpUsername    string `json:"smtp_username"`
	SmtpPassword    string `json:"smtp_passwd"`
}

type Configuration struct {
	ListenPort       int            `json:"listen_port"`
	ConnectionString string         `json:"rabbit_connection_string"`
	RetriesCount     int            `json:"retries_count"`
	Logstash         LogstashConfig `json:"logstash"`
	Email            EmailConfig    `json:"smtp"`
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
