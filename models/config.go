package models

import (
	"encoding/json"
	"os"
)

type Config struct {
	Geckoboard Geckoboard `json:"geckoboard"`
	Zendesk    Zendesk    `json:"zendesk"`
}

type Geckoboard struct {
	APIKey string `json:"api_key"`
	URL    string `json:"url"`
}

type Zendesk struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	APIKey   string `json:"api_key"`
	URL      string `json:"url"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
