package models

import (
	"testing"
)

type TestCase struct {
	DF       DateFilter
	Valid    bool
	ErrorMsg string
}

func TestValidateDateFilter(t *testing.T) {
	testCases := []TestCase{
		{
			DF: DateFilter{
				Attribute: "created",
				Unit:      "day",
				Past:      30,
			},
			Valid: true,
		},
		{
			DF: DateFilter{
				Attribute: "updated",
				Unit:      "month",
				Past:      1,
			},
			Valid: true,
		},
		{
			DF: DateFilter{
				Attribute: "solved",
				Unit:      "year",
				Past:      2,
			},
			Valid: true,
		},
		{
			DF: DateFilter{
				Attribute: "due_date",
				Unit:      "month",
				Past:      3,
			},
			Valid: true,
		},
		{
			DF: DateFilter{
				Attribute: "created",
				Custom:    ">2012-03-01",
			},
			Valid: true,
		},
		{
			DF: DateFilter{
				Attribute: "created",
				Custom:    "<2012-03-01",
			},
			Valid: true,
		},
		{
			DF: DateFilter{
				Attribute: "created",
				Custom:    ":2012-03-01",
			},
			Valid: true,
		},
		{
			DF: DateFilter{
				Unit: "month",
				Past: 2,
			},
			Valid: true,
		},
		{
			DF: DateFilter{
				Custom: ":2012-03-01",
			},
			Valid: true,
		},
		{
			DF:       DateFilter{},
			Valid:    false,
			ErrorMsg: "Date range or custom input is required",
		},
		{
			DF: DateFilter{
				Unit: "day",
			},
			Valid:    false,
			ErrorMsg: "Past is required to determine how many units to go into the past from today",
		},
		{
			DF: DateFilter{
				Attribute: "due_date",
				Past:      10,
			},
			Valid:    false,
			ErrorMsg: "Unit is required one of [minute hour day month year]",
		},
		{
			DF: DateFilter{
				Attribute: "due_date",
				Unit:      "dd",
				Past:      10,
			},
			Valid:    false,
			ErrorMsg: "Unit is required one of [minute hour day month year]",
		},
		{
			DF: DateFilter{
				Attribute: "duedate",
				Unit:      "day",
				Past:      10,
			},
			Valid:    false,
			ErrorMsg: "Attribute is required one of [created updated solved due_date]",
		},
		{
			DF: DateFilter{
				Custom: "2016-01-01",
			},
			Valid:    false,
			ErrorMsg: "Custom input requires the operator one of [< : >]",
		},
		{
			DF: DateFilter{
				Attribute: "created",
				Unit:      "day",
				Past:      10,
				Custom:    "2016-01-01",
			},
			Valid:    false,
			ErrorMsg: "Can't use both the unit, past and custom either unit & past or custom on its own",
		},
	}

	for _, tc := range testCases {
		err := tc.DF.Validate()

		if tc.Valid {
			if err != nil {
				t.Errorf("Expected no error but got %s", err)
			}
		} else {
			if err.Error() != tc.ErrorMsg {
				t.Errorf("Expected error %s but got %s", tc.ErrorMsg, err.Error())
			}
		}
	}
}

func TestDefaultAttributeWhenAttributeSupplied(t *testing.T) {
	df := DateFilter{
		Attribute: "solved",
		Unit:      "day",
		Past:      10,
		Custom:    "2016-01-01",
	}

	df.defaultAttribute()

	if df.Attribute != "solved" ||
		df.Unit != "day" ||
		df.Past != 10 ||
		df.Custom != "2016-01-01" {
		t.Errorf("Expected the date filter instance not to change but got %v", df)
	}
}

func TestDefaultAttribureWhenAttributeMissing(t *testing.T) {
	df := DateFilter{
		Unit:   "day",
		Past:   10,
		Custom: "2016-01-01",
	}

	df.defaultAttribute()

	if df.Attribute != "created" ||
		df.Unit != "day" ||
		df.Past != 10 ||
		df.Custom != "2016-01-01" {
		t.Errorf("Expected the attribute to change to created but got %v", df)
	}
}
