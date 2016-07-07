package zendesk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/geckoboard/zendesk_dataset/conf"
)

// Client contains the zendesk auth and where the client should paginate the results
type Client struct {
	Auth            conf.Auth
	PaginateResults bool
}

const basePath = "/api/v2"
const searchPath = "/search.json"

var scheme = "https"
var host = "%s.zendesk.com"
var client = &http.Client{Timeout: time.Second * 10}

// NewClient is a method to quickly generate a new client with just two args
func NewClient(auth *conf.Auth, paginateResults bool) *Client {
	return &Client{
		Auth:            *auth,
		PaginateResults: paginateResults,
	}
}

func (c *Client) buildRequest(method, path, queryParams, fullURL string) (*http.Request, error) {
	var uri string

	if fullURL != "" {
		uri = fullURL
	} else {
		u := url.URL{
			Scheme: scheme,
			Host:   fmt.Sprintf(host, c.Auth.Subdomain),
			Path:   basePath + path,
		}

		if queryParams != "" {
			q, err := url.ParseQuery(queryParams)
			if err != nil {
				return nil, err
			}

			u.RawQuery = q.Encode()
		}

		uri = u.String()
	}

	req, err := http.NewRequest(method, uri, nil)

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

// SearchTickets takes a string of the queryparams which consists of the zendesk
// query and returns a TicketPayload. If the Client specifies that it should
// paginate the results then this will utilize Zendesk's next_page attribute
// in the ticket payload until it returns at empty string.
func (c *Client) SearchTickets(queryParams string) (*TicketPayload, error) {
	res := TicketPayload{}
	var tp []Ticket

	if queryParams != "" {
		queryParams = "query=" + queryParams
	}

	totalCount, err := c.searchTickets(queryParams, "", &tp)
	if err != nil {
		return nil, err
	}

	res.Tickets = tp
	if c.PaginateResults {
		res.Count = len(tp)
	} else {
		res.Count = totalCount
	}

	return &res, nil
}

func (c *Client) searchTickets(queryParam, fullURL string, t *[]Ticket) (int, error) {
	req, err := c.buildRequest("GET", searchPath, queryParam, fullURL)

	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var tp TicketPayload
	err = json.NewDecoder(resp.Body).Decode(&tp)
	if err != nil {
		return 0, err
	}

	for _, tck := range tp.Tickets {
		*t = append(*t, tck)
	}

	if c.PaginateResults && tp.NextPage != "" {
		_, err = c.searchTickets("", tp.NextPage, t)
		if err != nil {
			return 0, err
		}
	}

	return tp.Count, nil
}
