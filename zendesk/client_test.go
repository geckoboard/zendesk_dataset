package zendesk

import (
	"testing"

	"github.com/jnormington/geckoboard_zendesk/conf"
)

type TestCase struct {
	Config      conf.Zendesk
	Method      string
	Path        string
	QueryParams string
	Expected    Expected
}

type Expected struct {
	Method     string
	FullPath   string
	AuthHeader string
}

func TestBuildRequest(t *testing.T) {
	testCases := []TestCase{
		{
			Config: conf.Zendesk{
				Auth: conf.Auth{
					Subdomain: "test",
					Email:     "test@example.com",
					APIKey:    "1234abc",
				},
			},
			Method:      "GET",
			Path:        searchPath,
			QueryParams: "type:ticket created<2017-01-01",
			Expected: Expected{
				Method:     "GET",
				FullPath:   "https://test.zendesk.com/api/v2/search.json?query=type%3Aticket+created%3C2017-01-01",
				AuthHeader: "Basic dGVzdEBleGFtcGxlLmNvbS90b2tlbjoxMjM0YWJj",
			},
		},
		{
			Config: conf.Zendesk{
				Auth: conf.Auth{
					Subdomain: "testdomain",
					Email:     "test@example.com",
					Password:  "9876cba",
				},
			},
			Method:      "POST",
			Path:        searchPath,
			QueryParams: "",
			Expected: Expected{
				Method:     "POST",
				FullPath:   "https://testdomain.zendesk.com/api/v2/search.json",
				AuthHeader: "Basic dGVzdEBleGFtcGxlLmNvbTo5ODc2Y2Jh",
			},
		},
	}

	for _, tc := range testCases {

		client := Client{
			Auth: tc.Config.Auth,
		}

		req, err := client.buildRequest(tc.Method, tc.Path, tc.QueryParams)

		if err != nil {
			t.Fatal(err)
		}

		if req.Method != tc.Expected.Method {
			t.Errorf("Expected request method to be %s but got %s", tc.Expected.Method, req.Method)
		}

		if req.URL.String() != tc.Expected.FullPath {
			t.Errorf("Expected url to be %s but got %s", tc.Expected.FullPath, req.URL.String())
		}

		actAuth := req.Header.Get("Authorization")

		if actAuth != tc.Expected.AuthHeader {
			t.Errorf("Expected basic auth header %s, but got %s", tc.Expected.AuthHeader, actAuth)
		}

	}
}
