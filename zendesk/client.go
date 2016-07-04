package zendesk

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jnormington/geckoboard_zendesk/conf"
)

type Client struct {
	Auth            conf.Auth
	PaginateResults bool
}

var domainURL = "https://%s.zendesk.com/api/v2"

const searchPath = "/search.json"

func (c *Client) buildRequest(method, path, queryParams string) (*http.Request, error) {
	uri := fmt.Sprintf(domainURL, c.Auth.Subdomain) + path

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
