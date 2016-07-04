package conf

import (
	"encoding/json"
	"testing"
	"time"
)

var staticTime = time.Date(2016, time.June, 01, 0, 0, 0, 0, time.FixedZone("BST", 1))

type TestCase struct {
	DF            DateFilter
	ExpectedQuery string
	Valid         bool
	ErrorMsg      string
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
			ErrorMsg: "Unit is required one of [day month year]",
		},
		{
			DF: DateFilter{
				Attribute: "due_date",
				Unit:      "dd",
				Past:      10,
			},
			Valid:    false,
			ErrorMsg: "Unit is required one of [day month year]",
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

func TestAddDateAPIFormat(t *testing.T) {
	testCases := []map[string]string{
		{"Input": `{"Past": 3, "Unit": "day"}`, "Output": "2016-05-29"},
		{"Input": `{"Past": 3, "Unit": "month"}`, "Output": "2016-03-01"},
		{"Input": `{"Past": 1, "Unit": "year"}`, "Output": "2015-06-01"},
		{"Input": `{"Past": 25, "Unit": "hour"}`, "Output": "2016-06-01"},
		{"Input": `{"Past": 2880, "Unit": "minute"}`, "Output": "2016-06-01"},
	}

	for _, tc := range testCases {
		var df DateFilter
		err := json.Unmarshal([]byte(tc["Input"]), &df)
		if err != nil {
			t.Fatal(err)
		}

		output := df.getDateAPIFormat(&staticTime)

		if output != tc["Output"] {
			t.Errorf("Expected output to be %s but got %s", tc["Output"], output)
		}
	}
}

func TestBuildQuery(t *testing.T) {
	testCases := []map[string]string{
		{"Input": `{"Past": 1, "Unit": "month"}`, "Output": "created>2016-05-01"},
		{"Input": `{"Attribute": "updated", "Past": 7, "Unit": "day"}`, "Output": "updated>2016-05-25"},
		{"Input": `{"Attribute": "solved", "Custom": ">2016-02-11"}`, "Output": "solved>2016-02-11"},
		{"Input": `{"Attribute": "due_date", "Custom": "<=2016-02-11"}`, "Output": "due_date<=2016-02-11"},
	}

	for _, tc := range testCases {
		var df DateFilter
		err := json.Unmarshal([]byte(tc["Input"]), &df)
		if err != nil {
			t.Fatal(err)
		}

		output := df.BuildQuery(&staticTime)

		if output != tc["Output"] {
			t.Errorf("Expected output to be %s but got %s", tc["Output"], output)
		}
	}
}

func TestDateRangeBuildQuery(t *testing.T) {
	dr1 := DateFilters{
		{
			Unit: month,
			Past: 2,
		},
		{
			Attribute: "updated",
			Custom:    "<2017-01-01",
		},
	}

	dr2 := DateFilters{
		{
			Unit: month,
			Past: 2,
		},
	}

	output1 := dr1.BuildQuery(&staticTime)
	if output1 != "created>2016-04-01 updated<2017-01-01" {
		t.Errorf("Built query output not matched got %s", output1)
	}

	output2 := dr2.BuildQuery(&staticTime)
	if output2 != "created>2016-04-01" {
		t.Errorf("Built query output not matched got %s", output2)
	}
}
