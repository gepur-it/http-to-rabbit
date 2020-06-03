package main

import (
	"os"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	log "github.com/sirupsen/logrus"
	"github.com/zbindenren/logrus_mail"
)

func initLogger(config Configuration) {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	// configure hooks
	hooks := make([]log.Hook, 0)

	if config.Logstash != (LogstashConfig{}) {
		log.Info("Adding logstash logging")
		hook, err := logrustash.NewHook(config.Logstash.Protocol, config.Logstash.Address, config.Logstash.ApplicationName)

		if err != nil {
			log.Error(err)
		} else {
			log.Info("Logstash logging added")
			hooks = append(hooks, hook)
		}
	}

	if config.Email != (EmailConfig{}) {
		log.Info("Adding email logging")

		hook, err := logrus_mail.NewMailAuthHook(
			config.Email.ApplicationName,
			config.Email.SmtpHost,
			config.Email.SmtpPort,
			config.Email.SmtpFrom,
			config.Email.SmtpTo,
			config.Email.SmtpUsername,
			config.Email.SmtpPassword,
		)

		if err != nil {
			log.Error(err)
		} else {
			log.Info("Email logging added")
			hooks = append(hooks, hook)
		}
	}

	for i := 0; i < len(hooks); i++ {
		log.AddHook(hooks[i])
	}
}
