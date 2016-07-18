package zendesk

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/geckoboard/zendesk_dataset/conf"
)

type TestCase struct {
	Config   conf.Auth
	Method   string
	FullURL  string
	Expected map[string]string
}

func TestBuildRequest(t *testing.T) {
	testCases := []TestCase{
		{
			Config: conf.Auth{
				Subdomain: "test",
				Email:     "test@example.com",
				APIKey:    "1234abc",
			},
			Method:  "GET",
			FullURL: "https://test.zendesk.com/api/v2/search.json?query=type%3Aticket+created%3C2017-01-01",
			Expected: map[string]string{
				"Method":     "GET",
				"FullPath":   "https://test.zendesk.com/api/v2/search.json?query=type%3Aticket+created%3C2017-01-01",
				"AuthHeader": "Basic dGVzdEBleGFtcGxlLmNvbS90b2tlbjoxMjM0YWJj",
			},
		},
		{
			Config: conf.Auth{
				Subdomain: "testdomain",
				Email:     "test@example.com",
				Password:  "9876cba",
			},
			Method:  "POST",
			FullURL: "https://testdomain.zendesk.com/api/v2/search.json",
			Expected: map[string]string{
				"Method":     "POST",
				"FullPath":   "https://testdomain.zendesk.com/api/v2/search.json",
				"AuthHeader": "Basic dGVzdEBleGFtcGxlLmNvbTo5ODc2Y2Jh",
			},
		},
		{
			Config: conf.Auth{
				Subdomain: "testdomain",
				Email:     "test@example.com",
				Password:  "9876cba",
			},
			Method:  "GET",
			FullURL: "http://blah.example.com?query=tags%3Atest",
			Expected: map[string]string{
				"Method":     "GET",
				"FullPath":   "http://blah.example.com?query=tags%3Atest",
				"AuthHeader": "Basic dGVzdEBleGFtcGxlLmNvbTo5ODc2Y2Jh",
			},
		},
	}

	for _, tc := range testCases {

		c := Client{
			Auth: tc.Config,
		}

		req, err := c.buildRequest(tc.Method, tc.FullURL)

		if err != nil {
			t.Fatal(err)
		}

		if req.Method != tc.Expected["Method"] {
			t.Errorf("Expected request method to be %s but got %s", tc.Expected["Method"], req.Method)
		}

		if req.URL.String() != tc.Expected["FullPath"] {
			t.Errorf("Expected url to be %s but got %s", tc.Expected["FullPath"], req.URL.String())
		}

		actAuth := req.Header.Get("Authorization")

		if actAuth != tc.Expected["AuthHeader"] {
			t.Errorf("Expected basic auth header %s, but got %s", tc.Expected["AuthHeader"], actAuth)
		}
	}
}

type STTestCase struct {
	Query                 string
	PaginateResults       bool
	RequestCount          int
	Requests              []Request
	ExpectedRequestCount  int
	ExpectedTicketPayload TicketPayload
	ExpectedTicketMetrics TicketMetrics
}

type Request struct {
	ReplaceBodyWithServer bool
	FullPath              string
	ResponseBody          string
}

var serverURL string

func TestSearchTickets(t *testing.T) {
	testCases := []STTestCase{
		{
			Query:                "type:ticket status:pending",
			ExpectedRequestCount: 1,
			Requests: []Request{
				{
					FullPath: "/api/v2/search.json?query=type%3Aticket+status%3Apending",
					//Mimic pagination not happening and the count being more than the first page
					ResponseBody: `{"results": [{"id": 1, "tags": ["beta", "test"]},` +
						`{"id": 2, "tags": ["expired", "test"]}], "count": 40, "next_page": ""}`,
				},
			},
			ExpectedTicketPayload: TicketPayload{
				Count: 40,
				Tickets: []Ticket{
					{ID: 1, Tags: []string{"beta", "test"}},
					{ID: 2, Tags: []string{"expired", "test"}},
				},
			},
		},
		{
			//Multi pages to return
			Query:                "type:ticket tags:important",
			PaginateResults:      true,
			ExpectedRequestCount: 2,
			Requests: []Request{
				{
					ReplaceBodyWithServer: true,
					FullPath:              "/api/v2/search.json?query=type%3Aticket+tags%3Aimportant",
					ResponseBody: `{"results": [{"id": 1, "tags": ["important", "test"]},` +
						`{"id": 2, "tags": ["important", "test"]}], "count": 2, "next_page": ` +
						`"%s/api/v2/search.json?page=2" }`,
				},
				{
					FullPath: "/api/v2/search.json?page=2",
					ResponseBody: `{"results": [{"id": 3, "tags": ["beta", "important"]},` +
						`{"id": 4, "tags": ["expired", "important"]}], "count": 2 }`,
				},
			},
			ExpectedTicketPayload: TicketPayload{
				Count: 4,
				Tickets: []Ticket{
					{ID: 1, Tags: []string{"important", "test"}},
					{ID: 2, Tags: []string{"important", "test"}},
					{ID: 3, Tags: []string{"beta", "important"}},
					{ID: 4, Tags: []string{"expired", "important"}},
				},
			},
		},
	}

	for _, tc := range testCases {
		server := buildServerWithExpectations(&tc, t)
		defer server.Close()

		//Reset the scheme and host back to original so not to break other tests
		defer func(h, s string) { host = h; scheme = s }(host, scheme)

		scheme = "http"
		host = "%s" + strings.Replace(server.URL, "http://", "", 1)
		serverURL = server.URL

		clt := Client{
			Auth: conf.Auth{
				Subdomain: "",
			},
			PaginateResults: tc.PaginateResults,
		}

		tp, err := clt.SearchTickets(&Query{Params: tc.Query})
		if err != nil {
			t.Fatal(err)
		}

		if tc.RequestCount != tc.ExpectedRequestCount {
			t.Errorf("Expected %d requests but got %d", tc.ExpectedRequestCount, tc.RequestCount)
		}

		if !reflect.DeepEqual(*tp, tc.ExpectedTicketPayload) {
			t.Errorf("Expected payload %v but got %v", tc.ExpectedTicketPayload, tp)
		}
	}
}

func TestTicketMetrics(t *testing.T) {
	splitTicketCount = 8
	testCases := []STTestCase{
		{
			Query:                "type:ticket",
			PaginateResults:      true,
			ExpectedRequestCount: 4,
			Requests: []Request{
				{
					ReplaceBodyWithServer: true,
					FullPath:              "/api/v2/search.json?query=type%3Aticket",
					ResponseBody: `{"results": [{"id": 1, "tags": ["important", "test"]},` +
						`{"id": 2, "tags": ["important", "test"]}], "count": 2, "next_page": "%s/api/v2/search.json?page=2" }`,
				},
				{
					FullPath: "/api/v2/search.json?page=2",
					ResponseBody: `{"results": [{"id": 3, "tags": ["beta", "important"]},
						{"id": 4, "tags": ["expired", "important"]},{"id":5},{"id":6},{"id":7},
						{"id":8},{"id":9},{"id":10},{"id":11},{"id":12},{"id":13},{"id":14}], "count": 20}`,
				},
				{
					ReplaceBodyWithServer: true,
					FullPath:              "/api/v2/tickets/show_many.json?ids=1%2C2%2C3%2C4%2C5%2C6%2C7%2C8%2C9&include=metric_sets",
					ResponseBody: `{"tickets":[{"metric_set": {"reply_time_in_minutes": {"calendar": 123},
									"full_resolution_time_in_minutes": {"business": 120, "calendar": 100}}}]}`,
				},
				{
					FullPath:     "/api/v2/tickets/show_many.json?ids=10%2C11%2C12%2C13%2C14&include=metric_sets",
					ResponseBody: `{"tickets":[{"metric_set": {"reply_time_in_minutes": {"calendar": 103}}}]}`,
				},
			},
			ExpectedTicketMetrics: TicketMetrics{
				Count: 2,
				Tickets: []Ticket{
					{
						Metrics: MetricSet{
							ReplyTime: SubTimeMetric{
								Calendar: 123,
							},
							FullResolutionTime: SubTimeMetric{
								Business: 120,
								Calendar: 100,
							},
						},
					},
					{
						Metrics: MetricSet{
							ReplyTime: SubTimeMetric{
								Calendar: 103,
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		server := buildServerWithExpectations(&tc, t)
		defer server.Close()

		//Reset the scheme and host back to original so not to break other tests
		defer func(h, s string) { host = h; scheme = s }(host, scheme)

		scheme = "http"
		host = "%s" + strings.Replace(server.URL, "http://", "", 1)
		serverURL = server.URL

		clt := Client{
			Auth: conf.Auth{
				Subdomain: "",
			},
			PaginateResults: tc.PaginateResults,
		}

		tp, err := clt.TicketMetrics(&Query{Params: tc.Query})
		if err != nil {
			t.Fatal(err)
		}

		if tc.RequestCount != tc.ExpectedRequestCount {
			t.Errorf("Expected %d requests but got %d", tc.ExpectedRequestCount, tc.RequestCount)
		}

		if !reflect.DeepEqual(*tp, tc.ExpectedTicketMetrics) {
			t.Errorf("Expected payload %#v but got %#v", tc.ExpectedTicketMetrics, tp)
		}
	}
}

func TestBuildURL(t *testing.T) {
	baseURL := "https://test.zendesk.com/api/v2"

	c := Client{
		Auth: conf.Auth{
			Subdomain: "test",
		},
	}

	tests := []struct {
		in          [2]string
		extraParams map[string]string
		out         string
		err         string
	}{
		{
			in:  [2]string{searchPath, "created>2016-01-01 tags:beta"},
			out: baseURL + "/search.json?query=created%3E2016-01-01+tags%3Abeta",
		},
		{
			in:  [2]string{"", "created>2016-01-01 tags:test"},
			err: "Endpoint is required to build url",
		},
		{
			in:  [2]string{"/ticket_metrics.json", ""},
			out: baseURL + "/ticket_metrics.json",
		},
		{
			in:          [2]string{"/tickets/show_many.json", ""},
			extraParams: map[string]string{"include": "metric_sets", "sort_by": "created_at"},
			out:         baseURL + "/tickets/show_many.json?include=metric_sets&sort_by=created_at",
		},
	}

	for _, tc := range tests {
		q := &Query{
			Endpoint:    tc.in[0],
			Params:      tc.in[1],
			ExtraParams: tc.extraParams,
		}
		out, err := c.buildURL(q)

		if tc.out != "" && tc.out != out {
			t.Errorf("Expected output to be %s but got %s", tc.out, out)
		}

		if tc.err == "" && err != nil {
			t.Errorf("Expected no errors but got %s", err.Error())
		}

		if tc.err != "" && tc.err != err.Error() {
			t.Errorf("Expected error %s but got %s", tc.err, err)
		}
	}
}

func buildServerWithExpectations(s *STTestCase, t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for _, e := range s.Requests {
			if r.URL.String() == e.FullPath {
				s.RequestCount++
				if e.ReplaceBodyWithServer {
					fmt.Fprintf(w, e.ResponseBody, serverURL)
				} else {
					fmt.Fprintf(w, e.ResponseBody)
				}
			}
		}
	}))

	return server
}
