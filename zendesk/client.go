package zendesk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/geckoboard/zendesk_dataset/conf"
)

type client struct {
	Auth            conf.Auth
	PaginateResults bool
}

type query struct {
	Endpoint string
	Params   string
}

const basePath = "/api/v2"
const searchPath = "/search.json"

var scheme = "https"
var host = "%s.zendesk.com"
var httpClt = &http.Client{Timeout: time.Second * 10}

func newClient(auth *conf.Auth, paginateResults bool) *client {
	return &client{
		Auth:            *auth,
		PaginateResults: paginateResults,
	}
}

func (c *client) buildURL(qy *query) (string, error) {
	if qy.Endpoint == "" {
		return "", errors.New("Endpoint is required to build url")
	}

	u := url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf(host, c.Auth.Subdomain),
		Path:   basePath + qy.Endpoint,
	}

	if qy.Params != "" {
		q, err := url.ParseQuery("query=" + qy.Params)
		if err != nil {
			return "", err
		}

		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

func (c *client) buildRequest(method, fullURL string) (*http.Request, error) {
	req, err := http.NewRequest(method, fullURL, nil)

	if err != nil {
		return nil, err
	}

	if c.Auth.Password != "" {
		req.SetBasicAuth(c.Auth.Email, c.Auth.Password)
	} else {
		req.SetBasicAuth(c.Auth.Email+"/token", c.Auth.APIKey)
	}

	return req, nil
}

func (c *client) searchTickets(q *query) (*ticketPayload, error) {
	var t []ticket

	q.Endpoint = searchPath
	var url, err = c.buildURL(q)
	if err != nil {
		return nil, err
	}

	for url != "" {
		req, err := c.buildRequest("GET", url)
		if err != nil {
			return nil, err
		}

		resp, err := httpClt.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		var tp ticketPayload
		err = json.NewDecoder(resp.Body).Decode(&tp)
		if err != nil {
			return nil, err
		}

		if c.PaginateResults {
			url = tp.NextPage
			t = append(t, tp.Tickets...)
		} else {
			return &tp, nil
		}
	}

	return &ticketPayload{Count: len(t), Tickets: t}, nil
}
