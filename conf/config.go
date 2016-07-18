package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config is the toplevel instance of the config containing
// both Geckoboard and Zendesk configurations.
type Config struct {
	Geckoboard Geckoboard `yaml:"geckoboard"`
	Zendesk    Zendesk    `yaml:"zendesk"`
}

// Geckoboard describes the authentication options.
type Geckoboard struct {
	APIKey string `yaml:"api_key"`
	URL    string `yaml:"url"`
}

// Auth makes up the Zendesk authentication options.
type Auth struct {
	Email     string `yaml:"email"`
	Password  string `yaml:"password"`
	APIKey    string `yaml:"api_key"`
	Subdomain string `yaml:"subdomain"`
}

// Zendesk contains Auth and a slice of Reports.
type Zendesk struct {
	Auth    Auth     `yaml:"auth"`
	Reports []Report `yaml:"reports"`
}

// Report describes the template to use and the filters to build for the Zendesk request.
type Report struct {
	Name          string       `yaml:"name"`
	DataSet       string       `yaml:"dataset"`
	GroupBy       GroupBy      `yaml:"group_by"`
	Filter        SearchFilter `yaml:"filter"`
	MetricOptions MetricOption `yaml:"metric_options"`
}

// GroupBy describes how a report should be grouped.
type GroupBy struct {
	Key  string `yaml:"key"`
	Name string `yaml:"name"`
}

// DisplayName returns the key if Name attribute is an empty string
// or returns the Name attribute specified by the user.
func (gb *GroupBy) DisplayName() string {
	if gb.Name == "" && gb.Key != "" {
		return gb.Key
	}

	return gb.Name
}

// LoadConfig take path and attempts to open the file and
// returns any errors that might occur with yaml syntax or file issues.
func LoadConfig(path string) (*Config, error) {
	var config Config

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
