package main

import (
	"flag"
	"log"

	"github.com/geckoboard/zendesk_dataset/conf"
	"github.com/geckoboard/zendesk_dataset/zendesk"
)

var configPath = flag.String("config", "./geckoboard_zendesk.conf", "Path to your geckoboard zendesk configuration")

func main() {
	flag.Parse()

	config, err := conf.LoadConfig(*configPath)

	if err != nil {
		log.Fatalf("ERRO: Problem with the config: %s\n", err)
	}

	if len(config.Zendesk.Reports) == 0 {
		log.Fatal("ERRO: You have no reports setup in your config under zendesk")
	}

	zendesk.HandleReports(config)
	log.Println("Completed processing all reports...")
}
