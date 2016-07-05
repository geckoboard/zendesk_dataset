package zendesk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jnormington/geckoboard_zendesk/conf"
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
					FullPath:     "/api/v2/search.json?query=type%3Aticket+created%3E2016-05-25",
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
	}

	for _, tc := range testCases {
		zserver := buildZendeskServerWithExpectations(&tc, t)
		gserver := buildGeckoboardServerWithExpectations(&tc, t)
		defer zserver.Close()
		defer gserver.Close()

		//Required by the client buildRequest method that %s
		domainURL = "%s" + zserver.URL + "/api/v2"
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
