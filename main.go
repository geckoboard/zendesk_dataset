package main

import (
	"flag"
	"log"

	"github.com/jnormington/geckoboard_zendesk/models"
)

var configPath = flag.String("config", "./geckoboard_zendesk.conf", "Path to your geckoboard zendesk configuration")

func main() {
	flag.Parse()

	_, err := models.LoadConfig(*configPath)

	if err != nil {
		log.Fatalf("Err: %s\n", err)
	}
}
