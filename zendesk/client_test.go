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
	Config      conf.Auth
	Method      string
	Path        string
	QueryParams string
	FullURL     string
	Expected    map[string]string
}

func TestBuildRequest(t *testing.T) {
	testCases := []TestCase{
		{
			Config: conf.Auth{
				Subdomain: "test",
				Email:     "test@example.com",
				APIKey:    "1234abc",
			},
			Method:      "GET",
			Path:        searchPath,
			QueryParams: "query=type:ticket created<2017-01-01",
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
			Method:      "POST",
			Path:        searchPath,
			QueryParams: "",
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
			Method:      "GET",
			Path:        searchPath,
			QueryParams: "",
			FullURL:     "http://blah.example.com?query=tags%3Atest",
			Expected: map[string]string{
				"Method":     "GET",
				"FullPath":   "http://blah.example.com?query=tags%3Atest",
				"AuthHeader": "Basic dGVzdEBleGFtcGxlLmNvbTo5ODc2Y2Jh",
			},
		},
	}

	for _, tc := range testCases {

		client := Client{
			Auth: tc.Config,
		}

		req, err := client.buildRequest(tc.Method, tc.Path, tc.QueryParams, tc.FullURL)

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

		scheme = "http"
		host = "%s" + strings.Replace(server.URL, "http://", "", 1)
		serverURL = server.URL

		clt := Client{
			Auth: conf.Auth{
				Subdomain: "",
			},
			PaginateResults: tc.PaginateResults,
		}

		tp, err := clt.SearchTickets(tc.Query)
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

func buildServerWithExpectations(s *STTestCase, t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.RequestCount++

		for _, e := range s.Requests {
			if r.URL.String() == e.FullPath {
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
