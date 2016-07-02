package zendesk

import (
	"testing"

	m "github.com/jnormington/geckoboard_zendesk/models"
)

type TestCase struct {
	Config   m.Zendesk
	Method   string
	Path     string
	Expected Expected
}

type Expected struct {
	Method     string
	Path       string
	AuthHeader string
}

func TestBuildRequestWithAuth(t *testing.T) {
	testCases := []TestCase{
		{
			Config: m.Zendesk{
				URL:    "https://test.domain.com",
				Email:  "test@example.com",
				APIKey: "1234abc",
			},
			Method: "GET",
			Path:   "/hello",
			Expected: Expected{
				Method:     "GET",
				Path:       "/hello",
				AuthHeader: "Basic dGVzdEBleGFtcGxlLmNvbS90b2tlbjoxMjM0YWJj",
			},
		},
		{
			Config: m.Zendesk{
				URL:      "https://test.domain.com",
				Email:    "test@example.com",
				Password: "9876cba",
			},
			Method: "POST",
			Path:   "/tickets/1.json",
			Expected: Expected{
				Method:     "POST",
				Path:       "/tickets/1.json",
				AuthHeader: "Basic dGVzdEBleGFtcGxlLmNvbTo5ODc2Y2Jh",
			},
		},
	}

	for _, tc := range testCases {

		req, err := buildRequestWithAuth(&tc.Config, tc.Method, tc.Config.URL+tc.Path)

		if err != nil {
			t.Fatal(err)
		}

		if req.Method != tc.Expected.Method {
			t.Errorf("Expected request method to be %s but got %s", tc.Expected.Method, req.Method)
		}

		expURL := tc.Config.URL + tc.Path

		if req.URL.String() != expURL {
			t.Errorf("Expected url to be %s but got %s", expURL, req.URL.String())
		}

		actAuth := req.Header.Get("Authorization")

		if actAuth != tc.Expected.AuthHeader {
			t.Errorf("Expected basic auth header %s, but got %s", tc.Expected.AuthHeader, actAuth)
		}

	}
}
