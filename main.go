package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/geckoboard/zendesk_dataset/conf"
	"github.com/geckoboard/zendesk_dataset/zendesk"
)

var (
	configPath     = flag.String("config", "./geckoboard_zendesk.conf", "Path to your geckoboard zendesk configuration")
	displayVersion = flag.Bool("version", false, "Prints version of Zendesk Dataset")
)

const version = "0.2.0"

func main() {
	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	config, err := conf.LoadConfig(*configPath)

	if err != nil {
		log.Fatalf("ERRO: Problem with the config: %s\n", err.Error())
	}

	if len(config.Zendesk.Reports) == 0 {
		log.Fatal("ERRO: You have no reports setup in your config under zendesk")
	}

	zendesk.HandleReports(config)
	log.Println("Completed processing all reports...")
}
