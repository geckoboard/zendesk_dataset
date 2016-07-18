package conf

import (
	"path"
	"reflect"
	"testing"
)

var configPath = "../fixtures"
var config = Config{
	Geckoboard: Geckoboard{
		APIKey: "Ap1K4y",
		URL:    "https://testing.geckoboardexample.com",
	},
	Zendesk: Zendesk{
		Auth: Auth{
			Email:     "test@example.com",
			APIKey:    "12345",
			Subdomain: "testing",
			Password:  "test",
		},
		Reports: []Report{
			{
				Name:    "report_1",
				DataSet: "dataset1",
				Filter: SearchFilter{
					Type: "ticket",
					DateRange: []DateFilter{
						{
							Attribute: created,
							Unit:      day,
							Past:      14,
						},
						{
							Attribute: created,
							Custom:    "<2017-01-01",
						},
					},
					Value: map[string]string{
						"status:": "open",
					},
					Values: map[string][]string{
						"tags:": []string{
							"beta",
							"freetrial",
							"account_expired",
						},
					},
				},
			},
			{
				Name:    "report_2",
				DataSet: "dataset2",
				Filter: SearchFilter{
					DateRange: []DateFilter{
						{
							Attribute: created,
							Unit:      day,
							Past:      14,
						},
					},
				},
				MetricOptions: MetricOption{
					Attribute: ReplyTime,
					Unit:      BusinessMetric,
					Grouping: []MetricGroup{
						{Unit: minute, From: 0, To: 1},
						{Unit: minute, From: 1, To: 8},
					},
				},
			},
		},
	},
}

func TestConfigLoadJSONFile(t *testing.T) {
	// Yaml lib can also load json just check its compatibility with our stuff still works.
	fileContents, err := LoadConfig(path.Join(configPath, "example.conf"))

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*fileContents, config) {
		t.Errorf("Expected json fileContents to be %v but got %v", config, fileContents)
	}
}

func TestConfigLoadYamlFile(t *testing.T) {
	fileContents, err := LoadConfig(path.Join(configPath, "example.yml"))

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*fileContents, config) {
		t.Errorf("Expected yaml fileContents to be %v but got %v", config, fileContents)
	}
}

func TestConfigLoadInvalidJson(t *testing.T) {
	_, err := LoadConfig(configPath + "example.bad")

	if err == nil {
		t.Errorf("Expected error but didn't get one")
	}
}

func TestConfigNotExists(t *testing.T) {
	_, err := LoadConfig("/tmp/somemissing.file")

	if err == nil {
		t.Errorf("Expected error but didn't get one")
	}
}
