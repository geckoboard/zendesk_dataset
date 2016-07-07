package zendesk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/geckoboard/zendesk_dataset/conf"
)

// Client contains the zendesk auth and where the client should paginate the results
type Client struct {
	Auth            conf.Auth
	PaginateResults bool
}

var domainURL = "https://%s.zendesk.com/api/v2"
var client = &http.Client{Timeout: time.Second * 10}

const searchPath = "/search.json"

// NewClient is a method to quickly generate a new client with just two args
func NewClient(auth *conf.Auth, paginateResults bool) *Client {
	return &Client{
		Auth:            *auth,
		PaginateResults: paginateResults,
	}
}

func (c *Client) buildRequest(method, path, queryParams string) (*http.Request, error) {
	var uri string
	domain := fmt.Sprintf(domainURL, c.Auth.Subdomain)

	if strings.Contains(path, domain) {
		uri = path
		path = ""
	} else if strings.Contains(queryParams, domain) {
		uri = queryParams
		queryParams = ""
	} else {
		uri = domain + path
	}

	if queryParams != "" {
		qp := "?query="
		qp += url.QueryEscape(queryParams)
		uri += qp
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

	totalCount, err := c.searchTickets(queryParams, &tp)
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

func (c *Client) searchTickets(queryParam string, t *[]Ticket) (int, error) {
	req, err := c.buildRequest("GET", searchPath, queryParam)

	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	var tp TicketPayload
	err = json.NewDecoder(resp.Body).Decode(&tp)
	if err != nil {
		return 0, err
	}

	for _, tck := range tp.Tickets {
		*t = append(*t, tck)
	}

	if c.PaginateResults && tp.NextPage != "" {
		_, err = c.searchTickets(tp.NextPage, t)
		if err != nil {
			return 0, err
		}
	}

	return tp.Count, nil
}
