package zendesk

import (
	"net/http"

	m "github.com/jnormington/geckoboard_zendesk/models"
)

func buildRequestWithAuth(conf *m.Zendesk, method, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, path, nil)

	if err != nil {
		return nil, err
	}

	if conf.Password != "" {
		req.SetBasicAuth(conf.Email, conf.Password)
	} else {
		req.SetBasicAuth(conf.Email+"/token", conf.APIKey)
	}

	return req, nil
}
