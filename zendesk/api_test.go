package zendesk

import (
	"testing"

	"github.com/geckoboard/zendesk_dataset/conf"
)

func TestSubMetricValue(t *testing.T) {
	testCases := []struct {
		ticket Ticket
		in     conf.MetricAttribute
		out    SubTimeMetric
	}{
		{
			ticket: Ticket{Metrics: MetricSet{ReplyTime: SubTimeMetric{Business: 1, Calendar: 11}}},
			in:     conf.ReplyTime,
			out:    SubTimeMetric{Business: 1, Calendar: 11},
		},
		{
			ticket: Ticket{Metrics: MetricSet{FirstResolutionTime: SubTimeMetric{Business: 2, Calendar: 22}}},
			in:     conf.FirstResolutionTime,
			out:    SubTimeMetric{Business: 2, Calendar: 22},
		},
		{
			ticket: Ticket{Metrics: MetricSet{FullResolutionTime: SubTimeMetric{Business: 3, Calendar: 33}}},
			in:     conf.FullResolutionTime,
			out:    SubTimeMetric{Business: 3, Calendar: 33},
		},
		{
			ticket: Ticket{Metrics: MetricSet{AgentWaitTime: SubTimeMetric{Business: 4, Calendar: 44}}},
			in:     conf.AgentWaitTime,
			out:    SubTimeMetric{Business: 4, Calendar: 44},
		},
		{
			ticket: Ticket{Metrics: MetricSet{RequesterWaitTime: SubTimeMetric{Business: 5, Calendar: 55}}},
			in:     conf.RequesterWaitTime,
			out:    SubTimeMetric{Business: 5, Calendar: 55},
		},
		{
			ticket: Ticket{Metrics: MetricSet{OnHoldTime: SubTimeMetric{Business: 6, Calendar: 66}}},
			in:     conf.OnHoldTime,
			out:    SubTimeMetric{Business: 6, Calendar: 66},
		},
	}

	for i, tc := range testCases {
		out := tc.ticket.subTimeMetric(tc.in)

		if out == nil {
			t.Fatalf("Expected %#v but got nil", tc.out)
		}

		if *out != tc.out {
			t.Errorf("[spec %d] Expected %#v but got %#v", i, tc.out, *out)
		}
	}
}
