package http

import (
	"net/http"

	"pkg.glorieux.io/mantra"
)

type httpClient struct {
	*http.Client
}

// NewClient returns a new HTTP client
func NewClient() mantra.Service {
	return &httpClient{
		Client: &http.Client{},
	}
}

func (*httpClient) Stop() error {
	return nil
}

func (*httpClient) String() string {
	return "HTTPClient"
}
