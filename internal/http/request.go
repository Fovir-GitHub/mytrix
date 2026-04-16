package http

import (
	"fmt"
	"io"
	"net/http"
)

type httpMethod struct{ v string }

var (
	MethodGet  = httpMethod{http.MethodGet}
	MethodPost = httpMethod{http.MethodPost}
)

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
