package conf

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

type Client struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	APIKey    string `json:"api_key"`
	Subdomain string `json:"subdomain"`
}

type Zendesk struct {
	Client  Client   `json:"auth"`
	Reports []Report `json:"reports"`
}

type Report struct {
	Name    string       `json:"name"`
	DataSet string       `json:"dataset"`
	Filter  SearchFilter `json:"filter"`
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
