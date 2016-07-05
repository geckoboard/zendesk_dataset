package main

import (
	"flag"
	"log"

	"github.com/jnormington/geckoboard_zendesk/conf"
	"github.com/jnormington/geckoboard_zendesk/zendesk"
)

var configPath = flag.String("config", "./geckoboard_zendesk.conf", "Path to your geckoboard zendesk configuration")

func main() {
	flag.Parse()

	config, err := conf.LoadConfig(*configPath)

	if err != nil {
		log.Fatalf("Err: Problem with the config: %s\n", err)
	}

	if len(config.Zendesk.Reports) == 0 {
		log.Fatal("Err: You have no reports setup in your config under zendesk")
	}

	err = zendesk.HandleReports(config)

	if err != nil {
		log.Fatalf("Err: %s\n", err)
	}

	log.Println("Completed with no errors, apparently...")
}
