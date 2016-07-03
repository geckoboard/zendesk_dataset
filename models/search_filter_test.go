package models

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
				DateRange: []DateFilter{},
				Value: map[string]string{
					"status": "solved",
				},
			},
			Valid:    false,
			ErrorMsg: "The key 'status' is missing an operator as the last character choose one of [> : < >= <=]",
		},
		{
			SF: SearchFilter{
				DateRange: []DateFilter{},
				Value: map[string]string{
					"status>": "",
				},
			},
			Valid:    false,
			ErrorMsg: "The key 'status>' has no value to check against",
		},
		{
			SF: SearchFilter{
				DateRange: []DateFilter{},
				Values: map[string][]string{
					"tags:": []string{},
				},
			},
			Valid:    false,
			ErrorMsg: "The key 'tags:' has no values to check against",
		},
		{
			SF: SearchFilter{
				DateRange: []DateFilter{},
				Value: map[string]string{
					"": "open",
				},
			},
			Valid:    false,
			ErrorMsg: "Missing key for value 'open'",
		},
		{
			SF: SearchFilter{
				DateRange: []DateFilter{},
				Values: map[string][]string{
					"tags": []string{"trial_expired", "emea"},
				},
			},
			Valid:    false,
			ErrorMsg: "The key 'tags' is missing an operator as the last character choose one of [> : < >= <=]",
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
