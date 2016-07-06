package conf

import "testing"

type SearchTestCase struct {
	SF       SearchFilter
	Valid    bool
	ErrorMsg string
}

func TestValidate(t *testing.T) {
	testCases := []SearchTestCase{
		{
			SF:    SearchFilter{},
			Valid: true,
		},
		{
			SF: SearchFilter{
				Type:      "ticket",
				DateRange: []DateFilter{},
			},
			Valid: true,
		},
		{
			SF: SearchFilter{
				DateRange: []DateFilter{},
				Value: map[string]string{
					"tags:": "Test",
				},
			},
			Valid: true,
		},
		{
			SF: SearchFilter{
				DateRange: []DateFilter{},
				Value: map[string]string{
					"status>": "pending",
				},
			},
			Valid: true,
		},
		{
			SF: SearchFilter{
				DateRange: []DateFilter{},
				Value: map[string]string{
					"status<": "solved",
				},
			},
			Valid: true,
		},
		{
			SF: SearchFilter{
				DateRange: []DateFilter{},
				Values: map[string][]string{
					"tags:": []string{"trial_expired", "emea"},
				},
			},
			Valid: true,
		},
		{
			SF: SearchFilter{
				DateRange: []DateFilter{},
				Value:     map[string]string{},
				Values:    map[string][]string{},
			},
			Valid: true,
		},
		{
			SF: SearchFilter{
				DateRange: []DateFilter{
					{
						Attribute: "created",
						Custom:    ">2016-01-01",
					},
					{
						Unit: "month",
					},
				},
			},
			Valid:    false,
			ErrorMsg: "Past is required to determine how many units to go into the past from today",
		},
	}

	for _, tc := range testCases {
		err := tc.SF.Validate()

		if tc.Valid {
			if err != nil {
				t.Errorf("Expected no error but got %s", err)
			}
		} else {
			if err == nil {
				t.Errorf("%v Expected error but got none", tc)
			} else if err.Error() != tc.ErrorMsg {
				t.Errorf("Expected error %s but got %s", tc.ErrorMsg, err)
			}
		}
	}
}

func TestUpdatesTypeToTicket(t *testing.T) {
	sf := SearchFilter{
		Type: "blah",
	}

	sf.defaultType()
	if sf.Type != "ticket" {
		t.Errorf("Expected the search filter to have been updated to ticket but wasn't")
	}

}

type BuildQueryTC struct {
	SF     SearchFilter
	Output string
}

func TestSearchFilterBuildQuery(t *testing.T) {
	testCases := []BuildQueryTC{
		{
			SF: SearchFilter{
				DateRange: DateFilters{},
			},
			Output: "type:ticket",
		},
		{
			SF: SearchFilter{
				DateRange: DateFilters{},
				Values: map[string][]string{
					"groups:": []string{
						"first_line",
						"second_line",
					},
				},
			},
			Output: "type:ticket groups:first_line groups:second_line",
		},
		{
			SF: SearchFilter{
				DateRange: DateFilters{
					{
						Unit: month,
						Past: 2,
					},
					{
						Attribute: "updated",
						Custom:    "<2017-01-01",
					},
				},
			},
			Output: "type:ticket created>=2016-04-01 updated<2017-01-01",
		},
		{
			SF: SearchFilter{
				DateRange: DateFilters{
					{
						Unit: month,
						Past: 2,
					},
					{
						Attribute: "updated",
						Custom:    "<2017-01-01",
					},
				},
				Value: map[string]string{
					"tags:": "beta",
				},
			},
			Output: "type:ticket created>=2016-04-01 updated<2017-01-01 tags:beta",
		},
		{
			SF: SearchFilter{
				Type: "user",
				DateRange: DateFilters{
					{
						Unit: month,
						Past: 2,
					},
					{
						Attribute: "updated",
						Custom:    "<2017-01-01",
					},
				},
				Value: map[string]string{
					"tags:":   "beta not",
					"status>": "pending",
				},
			},
			Output: `type:user created>=2016-04-01 updated<2017-01-01 status>pending tags:"beta not"`,
		},
		{
			SF: SearchFilter{
				DateRange: DateFilters{
					{
						Unit: month,
						Past: 1,
					},
				},
				Values: map[string][]string{
					"tags:": []string{
						"expired",
						"freetrial",
						"test_user",
					},
				},
			},
			Output: "type:ticket created>=2016-05-01 tags:expired tags:freetrial tags:test_user",
		},
		{
			SF: SearchFilter{
				DateRange: DateFilters{
					{
						Attribute: "updated",
						Custom:    "<2017-01-01",
					},
				},
				Values: map[string][]string{
					"tags:": []string{
						"expired",
					},
					"status:": []string{
						"pending",
						"solved",
					},
				},
			},
			Output: "type:ticket updated<2017-01-01 status:pending status:solved tags:expired",
		},
		{
			SF: SearchFilter{
				DateRange: DateFilters{
					{
						Attribute: "updated",
						Custom:    "<2017-01-01",
					},
				},
				Value: map[string]string{
					"tags:":   "beta urgent",
					"groups>": "firstline",
				},
				Values: map[string][]string{
					"tags:": []string{
						"pending solved",
						"planone",
					},
				},
			},
			Output: `type:ticket updated<2017-01-01 groups>firstline tags:"beta urgent" tags:"pending solved" tags:planone`,
		},
	}

	for _, tc := range testCases {
		output := tc.SF.BuildQuery(&staticTime)

		if output != tc.Output {
			t.Errorf("Expected output %s but got %s", tc.Output, output)
		}
	}
}
