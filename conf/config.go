package conf

import (
	"encoding/json"
	"os"
)

// Config is the toplevel instance of the config containing
// both Geckoboard and Zendesk configurations
type Config struct {
	Geckoboard Geckoboard `json:"geckoboard"`
	Zendesk    Zendesk    `json:"zendesk"`
}

// Geckoboard describes the authentication options
type Geckoboard struct {
	APIKey string `json:"api_key"`
	URL    string `json:"url"`
}

// Auth makes up the Zendesk authentication options
type Auth struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	APIKey    string `json:"api_key"`
	Subdomain string `json:"subdomain"`
}

// Zendesk contains Auth and a slice of Reports
type Zendesk struct {
	Auth    Auth     `json:"auth"`
	Reports []Report `json:"reports"`
}

// Report describes the template to use and the filters to build for the Zendesk request.
type Report struct {
	Name    string       `json:"name"`
	DataSet string       `json:"dataset"`
	GroupBy GroupBy      `json:"group_by"`
	Filter  SearchFilter `json:"filter"`
}

// GroupBy describes how a report should be grouped
type GroupBy struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// DisplayName returns the key if Name attribute is an empty string
// or returns the Name attribute specified by the user
func (gb *GroupBy) DisplayName() string {
	if gb.Name == "" && gb.Key != "" {
		return gb.Key
	}

	return gb.Name
}

// LoadConfig take path and attempts to open the file and
// returns any errors that might occur with json syntax or file issues
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
