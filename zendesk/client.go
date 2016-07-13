package zendesk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
	Endpoint    string
	Params      string
	ExtraParams map[string]string
}

const (
	basePath    = "/api/v2"
	searchPath  = "/search.json"
	ticketsPath = "/tickets/show_many.json"
)

var (
	scheme  = "https"
	host    = "%s.zendesk.com"
	httpClt = &http.Client{Timeout: time.Second * 10}
)

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

	q := url.Values{}
	var err error
	if qy.Params != "" {
		q, err = url.ParseQuery("query=" + qy.Params)
		if err != nil {
			return "", err
		}
	}

	// Add any addtional params that don't require query=.
	if len(qy.ExtraParams) > 0 {
		for k, v := range qy.ExtraParams {
			q.Add(k, v)
		}
	}

	u.RawQuery = q.Encode()

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

// TicketMetrics takes a query and returns TicketMetrics or an error if it
// occurs. The ticket metrics utilises two endpoints first it uses SearchTickets
// to get all the ticket IDs and then makes a request on the tickets/show_many.json
// sideloading the metric sets. This allows greater flexibility on the filters
// possible to get the metrics you require specifically based on a SearchFilter
func (c *Client) TicketMetrics(q *Query) (*TicketMetrics, error) {
	//Use search API to filter tickets
	tp, err := c.SearchTickets(q)
	if err != nil {
		return nil, err
	}

	//Extract all the ticket ids
	var bf bytes.Buffer
	for i, t := range tp.Tickets {
		bf.WriteString(strconv.Itoa(t.ID))
		if i != len(tp.Tickets)-1 {
			bf.WriteString(",")
		}
	}

	qy := &Query{
		Endpoint: ticketsPath,
		ExtraParams: map[string]string{
			"include": "metric_sets",
			"ids":     bf.String(),
		},
	}

	url, err := c.buildURL(qy)
	var tickets []Ticket
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

		var tm TicketMetrics
		err = json.NewDecoder(resp.Body).Decode(&tm)
		if err != nil {
			return nil, err
		}

		url = tm.NextPage
		tickets = append(tickets, tm.Tickets...)
	}

	return &TicketMetrics{Count: len(tickets), Tickets: tickets}, nil
}
