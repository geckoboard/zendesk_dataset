package zendesk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/geckoboard/zendesk_dataset/conf"
)

type ReportTestCase struct {
	RequestCount              int
	ZendeskRequests           []ERequest
	GeckoboardRequests        []ERequest
	ExpectedTotalRequestCount int
	Config                    conf.Config
}

type ERequest struct {
	FullPath     string
	RequestBody  string
	ResponseBody string
}

func TestHandleReports(t *testing.T) {
	testCases := []ReportTestCase{
		{
			ExpectedTotalRequestCount: 4,
			ZendeskRequests: []ERequest{
				{
					FullPath: "/api/v2/search.json?query=type%3Aticket+tags%3Abeta",
					ResponseBody: `{"results": [{"id": 1, "tags": ["beta", "test"]},` +
						`{"id": 2, "tags": ["expired", "beta"]}], "count": 2, "next_page":"/api/v2/never_called"}`,
				},
				{
					FullPath:     "/api/v2/search.json?query=type%3Aticket+tags%3Atest",
					ResponseBody: `{"results": [{"id": 3, "tags": ["test"]}], "count": 1}`,
				},
			},
			GeckoboardRequests: []ERequest{
				{
					FullPath: "/datasets/my.dataset_report",
					RequestBody: `{"id":"my.dataset_report","fields":{"grouped_by":{"name":"Tags","type":"string"},` +
						`"ticket_count":{"name":"Ticket Count","type":"number"}},` +
						`"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
					ResponseBody: "{}\n",
				},
				{
					FullPath:     "/datasets/my.dataset_report/data",
					RequestBody:  `{"data":[{"grouped_by":"beta","ticket_count":2},{"grouped_by":"test","ticket_count":1}]}`,
					ResponseBody: "{}\n",
				},
			},
			Config: conf.Config{
				Geckoboard: conf.Geckoboard{
					URL: "",
				},
				Zendesk: conf.Zendesk{
					Reports: []conf.Report{
						{
							Name:    "ticket_counts",
							DataSet: "my.dataset_report",
							GroupBy: conf.GroupBy{
								Key:  "tags:",
								Name: "Tags",
							},
							Filter: conf.SearchFilter{
								Values: map[string][]string{
									"tags:": []string{"beta", "test"},
								},
							},
						},
					},
				},
			},
		},
		{
			ExpectedTotalRequestCount: 3,
			ZendeskRequests: []ERequest{
				{
					FullPath:     "/api/v2/search.json?query=type%3Aticket+created%3E%3D2016-05-25",
					ResponseBody: `{"results": [{"id": 1},{"id":2},{"id":3}, {"id": 4}],"count": 4, "next_page":"/api/v2/never_called"}`,
				},
			},
			GeckoboardRequests: []ERequest{
				{
					FullPath: "/datasets/tickets_in_last_7_days",
					RequestBody: `{"id":"tickets_in_last_7_days","fields":{"grouped_by":{"name":"All","type":"string"},` +
						`"ticket_count":{"name":"Ticket Count","type":"number"}},` +
						`"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
					ResponseBody: "{}\n",
				},
				{
					FullPath:     "/datasets/tickets_in_last_7_days/data",
					RequestBody:  `{"data":[{"grouped_by":"All","ticket_count":4}]}`,
					ResponseBody: "{}\n",
				},
			},
			Config: conf.Config{
				Geckoboard: conf.Geckoboard{
					URL: "",
				},
				Zendesk: conf.Zendesk{
					Reports: []conf.Report{
						{
							Name:    "ticket_counts",
							DataSet: "tickets_in_last_7_days",
							Filter: conf.SearchFilter{
								DateRange: conf.DateFilters{
									{
										Unit: "day",
										Past: 7,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			ExpectedTotalRequestCount: 4,
			ZendeskRequests: []ERequest{
				{
					FullPath:     "/api/v2/search.json?query=type%3Aticket+created%3E%3D2016-05-29",
					ResponseBody: `{"results":[ {"id": 1},{"id": 2},{"id": 3},{"id": 4}]}`,
				},
				{
					FullPath: "/api/v2/tickets/show_many.json?ids=1%2C2%2C3%2C4&include=metric_sets",
					ResponseBody: `{"tickets":[
					{"created_at": "2016-06-29T19:59:14Z", "metric_set": {"reply_time_in_minutes": {"calendar": 70, "business": 59}}},
					{"created_at": "2016-06-29T19:59:14Z", "metric_set": {"reply_time_in_minutes": {"calendar": 120, "business": 60}}},
					{"created_at": "2016-07-30T19:59:14Z", "metric_set": {"reply_time_in_minutes": {"calendar": 181, "business": 121}}},
					{"created_at": "2016-07-30T19:59:14Z", "metric_set": {"reply_time_in_minutes": {"calendar": 185, "business": 480}}}
					] }`,
				},
			},
			GeckoboardRequests: []ERequest{
				{
					FullPath: "/datasets/ticket_metrics_in_last_3days",
					RequestBody: `{"id":"ticket_metrics_in_last_3days","fields":{"count":{"name":"Count","type":"number"},` +
						`"grouping":{"name":"Grouping","type":"string"}},"created_at":"0001-01-01T00:00:00Z",` +
						`"updated_at":"0001-01-01T00:00:00Z"}`,
					ResponseBody: "{}\n",
				},
				{
					FullPath:     "/datasets/ticket_metrics_in_last_3days/data",
					RequestBody:  `{"data":[{"grouping":"0-1 hour","count":1},{"grouping":"1-8 hours","count":2}]}`,
					ResponseBody: "{}\n",
				},
			},
			Config: conf.Config{
				Geckoboard: conf.Geckoboard{
					URL: "",
				},
				Zendesk: conf.Zendesk{
					Reports: []conf.Report{
						{
							Name:    "detailed_metrics",
							DataSet: "ticket_metrics_in_last_3days",
							Filter: conf.SearchFilter{
								DateRange: conf.DateFilters{
									{
										Unit: "day",
										Past: 3,
									},
								},
							},
							MetricOptions: conf.MetricOption{
								Attribute: conf.ReplyTime,
								Unit:      conf.BusinessMetric,
								Grouping: []conf.MetricGroup{
									{Unit: "hour", From: 0, To: 1},
									{Unit: "hour", From: 1, To: 8},
								},
							},
						},
					},
				},
			},
		},
		{
			ExpectedTotalRequestCount: 4,
			ZendeskRequests: []ERequest{
				{
					FullPath:     "/api/v2/search.json?query=type%3Aticket+created%3E%3D2016-05-29",
					ResponseBody: `{"results":[ {"id": 1},{"id": 2},{"id": 3},{"id": 4}]}`,
				},
				{
					FullPath: "/api/v2/tickets/show_many.json?ids=1%2C2%2C3%2C4&include=metric_sets",
					ResponseBody: `{"tickets":[
					{"created_at": "2016-06-29T19:59:14Z", "metric_set": {"first_resolution_time_in_minutes": {"calendar": 70, "business": 59}}},
					{"created_at": "2016-06-29T19:59:14Z", "metric_set": {"first_resolution_time_in_minutes": {"calendar": 120, "business": 60}}},
					{"created_at": "2016-07-30T19:59:14Z", "metric_set": {"first_resolution_time_in_minutes": {"calendar": 181, "business": 121}}},
					{"created_at": "2016-07-30T19:59:14Z", "metric_set": {"first_resolution_time_in_minutes": {"calendar": 480, "business": 480}}}
					] }`,
				},
			},
			GeckoboardRequests: []ERequest{
				{
					FullPath: "/datasets/ticket_metrics_in_last_3days",
					RequestBody: `{"id":"ticket_metrics_in_last_3days","fields":{"count":{"name":"Count","type":"number"},` +
						`"grouping":{"name":"Grouping","type":"string"}},"created_at":"0001-01-01T00:00:00Z",` +
						`"updated_at":"0001-01-01T00:00:00Z"}`,
					ResponseBody: "{}\n",
				},
				{
					FullPath:     "/datasets/ticket_metrics_in_last_3days/data",
					RequestBody:  `{"data":[{"grouping":"0-1 hour","count":0},{"grouping":"60-480 minutes","count":3},{"grouping":"480-800 minutes","count":1}]}`,
					ResponseBody: "{}\n",
				},
			},
			Config: conf.Config{
				Geckoboard: conf.Geckoboard{
					URL: "",
				},
				Zendesk: conf.Zendesk{
					Reports: []conf.Report{
						{
							Name:    "detailed_metrics",
							DataSet: "ticket_metrics_in_last_3days",
							Filter: conf.SearchFilter{
								DateRange: conf.DateFilters{
									{
										Unit: "day",
										Past: 3,
									},
								},
							},
							MetricOptions: conf.MetricOption{
								Attribute: conf.FirstResolutionTime,
								Unit:      conf.CalendarMetric,
								Grouping: []conf.MetricGroup{
									{Unit: "hour", From: 0, To: 1},
									{Unit: "minute", From: 60, To: 480},
									{Unit: "minute", From: 480, To: 800},
								},
							},
						},
					},
				},
			},
		},
		{
			ExpectedTotalRequestCount: 3,
			ZendeskRequests: []ERequest{
				{
					FullPath: "/api/v2/search.json?query=type%3Aticket+created%3E%3D2016-05-01",
					ResponseBody: `{"results": [{"created_at": "2016-06-29T19:59:14Z"},{"created_at": "2016-06-29T19:59:14Z"},{"created_at": "2016-06-30T19:59:14Z"},
					{"created_at": "2016-07-01T19:59:14Z"},{"created_at": "2016-07-01T19:59:14Z"},{"created_at": "2016-07-01T19:59:14Z"},
					{"created_at": "2016-07-01T19:59:14Z"},{"created_at": "2016-07-05T19:59:14Z"},{"created_at": "2016-07-04T19:59:14Z"}]}`,
				},
			},
			GeckoboardRequests: []ERequest{
				{
					FullPath: "/datasets/tickets.last.month.by.day",
					RequestBody: `{"id":"tickets.last.month.by.day","fields":{"count":{"name":"Ticket Count","type":"number"},"date":{"name":"Date","type":"date"}},` +
						`"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
					ResponseBody: "{}\n",
				},
				{
					FullPath:     "/datasets/tickets.last.month.by.day/data",
					RequestBody:  `{"data":[{"date":"2016-06-29","count":2},{"date":"2016-06-30","count":1},{"date":"2016-07-01","count":4},{"date":"2016-07-05","count":1},{"date":"2016-07-04","count":1}]}`,
					ResponseBody: "{}\n",
				},
			},
			Config: conf.Config{
				Geckoboard: conf.Geckoboard{
					URL: "",
				},
				Zendesk: conf.Zendesk{
					Reports: []conf.Report{
						{
							Name:    "ticket_counts_by_day",
							DataSet: "tickets.last.month.by.day",
							Filter: conf.SearchFilter{
								DateRange: conf.DateFilters{
									{
										Unit: "month",
										Past: 1,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		zserver := buildZendeskServerWithExpectations(&tc, t)
		gserver := buildGeckoboardServerWithExpectations(&tc, t)
		defer zserver.Close()
		defer gserver.Close()

		//Required by the client buildRequest method that %s
		scheme = "http"
		host = "%s" + strings.Replace(zserver.URL, "http://", "", 1)
		timeNow = time.Date(2016, 06, 01, 0, 0, 0, 0, time.UTC)
		tc.Config.Geckoboard.URL = gserver.URL

		HandleReports(&tc.Config)

		if tc.RequestCount != tc.ExpectedTotalRequestCount {
			t.Errorf("Expected %d requests but got %d", tc.ExpectedTotalRequestCount, tc.RequestCount)
		}
	}
}

func buildZendeskServerWithExpectations(s *ReportTestCase, t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for _, e := range s.ZendeskRequests {
			if r.URL.String() == e.FullPath {
				s.RequestCount++
				fmt.Fprintf(w, e.ResponseBody)
				return
			}
		}

		t.Errorf("Unexpected url recieved %s", r.URL)

		fmt.Fprintf(w, "{}\n")
	}))

	return server
}

func buildGeckoboardServerWithExpectations(s *ReportTestCase, t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		defer r.Body.Close()
		body := strings.TrimRight(string(b), "\n")

		for _, e := range s.GeckoboardRequests {
			if r.URL.String() == e.FullPath {
				s.RequestCount++

				if body != e.RequestBody {
					t.Errorf("Expected body %s but got %s", e.RequestBody, body)
				}

				fmt.Fprintf(w, e.ResponseBody)
				return
			}
		}

		t.Fatal("No matching requests found for: %v", r)
	}))

	return server
}
