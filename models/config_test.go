package models

import (
	"path"
	"reflect"
	"testing"
)

var configPath = "../fixtures"

func TestConfigLoadFromFile(t *testing.T) {
	config := Config{
		Geckoboard: Geckoboard{
			APIKey: "Ap1K4y",
			URL:    "https://testing.geckoboardexample.com",
		},
		Zendesk: Zendesk{
			Email:    "test@example.com",
			APIKey:   "12345",
			URL:      "http://testing.zendesk.com",
			Password: "test",
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
			},
		},
	}

	fileContents, err := LoadConfig(path.Join(configPath, "example.conf"))

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*fileContents, config) {
		t.Errorf("Expected fileContents to be %#v, but got %#v", config, fileContents)
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
