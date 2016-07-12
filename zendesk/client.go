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

// Client holds the Zendesk auth and whether the client should paginate.
type Client struct {
	Auth            conf.Auth
	PaginateResults bool
}

// Query holds the params and endpoint for which the buildURL method uses.
type Query struct {
	Endpoint string
	Params   string
}

const basePath = "/api/v2"
const searchPath = "/search.json"

var scheme = "https"
var host = "%s.zendesk.com"
var httpClt = &http.Client{Timeout: time.Second * 10}

func newClient(auth *conf.Auth, paginateResults bool) *Client {
	return &Client{
		Auth:            *auth,
		PaginateResults: paginateResults,
	}
}

func (c *Client) buildURL(qy *Query) (string, error) {
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

func (c *Client) buildRequest(method, fullURL string) (*http.Request, error) {
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

// SearchTickets takes a query object and returns a TicketPayload. If the Client
// specifies that it should paginate the results then it will utilize
// next_page attribute in the ticket payload until it returns at empty string with all the tickets.
// When not paginated it will return only the TicketPayload with the count
func (c *Client) SearchTickets(q *Query) (*TicketPayload, error) {
	var t []Ticket

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

		var tp TicketPayload
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

	return &TicketPayload{Count: len(t), Tickets: t}, nil
}
