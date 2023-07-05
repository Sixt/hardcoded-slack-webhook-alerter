package main

import (
	"log"

	"github.com/sixt/hardcoded-slack-webhook-alerter/client"
	"github.com/sixt/hardcoded-slack-webhook-alerter/config"
	"github.com/sixt/hardcoded-slack-webhook-alerter/scanner"
)

func main() {
	conf, err := config.LoadConfig(config.ConfigPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	//collect the results
	scanner := scanner.New(conf)
	scanner.Scan()

	client := client.NewSlackClient(conf.DryRun, conf.Message)

	//for each result, send a request to the webhook
	for _, res := range scanner.Results {
		if err := client.SendMessage(res); err != nil {
			log.Println(err.Error())
		}
	}
}
