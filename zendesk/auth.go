package zendesk

import (
	"net/http"

	"github.com/jnormington/geckoboard_zendesk/conf"
)

func buildRequestWithAuth(cf *conf.Zendesk, method, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, path, nil)

	if err != nil {
		return nil, err
	}

	if cf.Password != "" {
		req.SetBasicAuth(cf.Email, cf.Password)
	} else {
		req.SetBasicAuth(cf.Email+"/token", cf.APIKey)
	}

	return req, nil
}
