package conf

import "testing"

func TestGetUnitInMinutes(t *testing.T) {
	testCases := []struct {
		unit  calendarUnit
		value int
		out   int
	}{
		{unit: hour, value: 1, out: 60},
		{unit: hour, value: 0, out: 0},
		{unit: minute, value: 3, out: 3}, // Not supported returns -1
		{unit: day, value: 1, out: -1},   // Not supported returns -1
		{unit: year, value: 1, out: -1},
	}

	for _, tc := range testCases {
		out := getUnitInMinutes(tc.unit, tc.value)

		if out != tc.out {
			t.Errorf("Expected %d minutes but got %d", tc.out, out)
		}
	}
}

func TestValid(t *testing.T) {
	testCases := []struct {
		m   MetricOption
		err error
	}{
		{m: MetricOption{}, err: errEmptyMetricOptions},
		{m: MetricOption{Attribute: "Test"}, err: errInvalidAttribute},
		{m: MetricOption{Attribute: ReplyTime, Unit: "Test"}, err: errInvalidUnit},
		{m: MetricOption{Attribute: ReplyTime, Unit: BusinessMetric}},
		{m: MetricOption{Attribute: FirstResolutionTime, Unit: CalendarMetric}},
		{m: MetricOption{Attribute: FullResolutionTime, Unit: CalendarMetric}},
		{m: MetricOption{Attribute: AgentWaitTime, Unit: BusinessMetric}},
		{m: MetricOption{Attribute: RequesterWaitTime, Unit: CalendarMetric}},
		{m: MetricOption{Attribute: OnHoldTime, Unit: BusinessMetric}},
	}

	for i, tc := range testCases {
		err := tc.m.Valid()

		if tc.err == nil && err != nil {
			t.Errorf("[spec %d] Unexpected error got %s", i, err)
		}

		if tc.err != nil && err != tc.err {
			t.Errorf("[spec %d] Expected %s but got %s", i, tc.err, err)
		}
	}
}

func TestGroupingsValid(t *testing.T) {
	testCases := []struct {
		m   MetricOption
		err error
	}{
		{m: MetricOption{Grouping: []MetricGroup{}}, err: errEmptyGrouping},
		{m: MetricOption{Grouping: []MetricGroup{{Unit: "tst", From: 0, To: 1}}}, err: errInvalidGroupUnit},
		{m: MetricOption{Grouping: []MetricGroup{{Unit: day, From: 1, To: 3}}}, err: errInvalidGroupUnit},
		{m: MetricOption{Grouping: []MetricGroup{{Unit: minute, From: 1, To: 0}}}, err: errFromGreaterThanTo},
		{m: MetricOption{Grouping: []MetricGroup{{Unit: minute, From: 3, To: 3}}}, err: errFromEqualToTo},
		{m: MetricOption{Grouping: []MetricGroup{{Unit: minute, From: 2, To: 3}, {Unit: day}}}, err: errInvalidGroupUnit},
		{m: MetricOption{Grouping: []MetricGroup{{Unit: minute, From: 0, To: 1}}}},
		{m: MetricOption{Grouping: []MetricGroup{{Unit: hour, From: 0, To: 1}}}},
	}

	for i, tc := range testCases {
		err := tc.m.GroupingValid()

		if tc.err == nil && err != nil {
			t.Errorf("Unexpected error got %s", err)
		}

		if tc.err != nil && err != tc.err {
			t.Errorf("[spec %d] Expected error %s but got %s", i, tc.err, err)
		}
	}
}

func TestDisplayName(t *testing.T) {
	testCases := []struct {
		mg  MetricGroup
		out string
	}{
		{mg: MetricGroup{Unit: hour, From: 0, To: 1}, out: "0-1 hour"},
		{mg: MetricGroup{Unit: hour, From: 1, To: 2}, out: "1-2 hours"},
		{mg: MetricGroup{Unit: hour, From: 1, To: 5}, out: "1-5 hours"},
		{mg: MetricGroup{Unit: minute, From: 0, To: 1}, out: "0-1 minute"},
		{mg: MetricGroup{Unit: minute, From: 0, To: 2}, out: "0-2 minutes"},
		{mg: MetricGroup{Unit: minute, From: 60, To: 120}, out: "60-120 minutes"},
	}

	for i, tc := range testCases {
		out := tc.mg.DisplayName()

		if tc.out != out {
			t.Errorf("[spec %d] Expected output %s but got %s", i, tc.out, out)
		}
	}
}
