package http

import (
	"fmt"
	"io"
	"net/http"
)

// httpMethod represents an HTTP request method type.
type httpMethod struct{ v string }

var (
	// MethodGet represents an HTTP GET request.
	MethodGet = httpMethod{http.MethodGet}
	// MethodPost represents an HTTP POST request.
	MethodPost = httpMethod{http.MethodPost}
)

// NewRequest creates a new HTTP request with the specified method, URL, body, and headers.
// It returns an error if the request cannot be created.
func NewRequest(method httpMethod, url string, body io.Reader, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method.v, url, body)
	if err != nil {
		return nil, fmt.Errorf("new http request failed (method=%s, url=%s): %w", method.v, url, err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}
