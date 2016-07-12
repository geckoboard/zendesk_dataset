package geckoboard

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendRequest(t *testing.T) {

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "bad basic syntax", http.StatusBadRequest)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if !(len(pair) == 2 && pair[0] == "secret" && pair[1] == "") {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))

	client := New(Config{URL: s.URL, Key: "secret"})
	resp, err := client.sendNewRequest("GET", "/foobar", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d status code, got %d", http.StatusOK, resp.StatusCode)
	}
}
