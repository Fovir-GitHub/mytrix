// Package http provides HTTP client functionality with logging and error handling.
// It wraps the standard library HTTP client with additional features like timeout configuration,
// request/response logging, and JSON decoding.
package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Fovir-GitHub/mytrix/internal/config"
)

// Client is a wrapper around *http.Client with timeout configuration and logging.
// It provides methods for making HTTP requests with automatic logging and error handling.
type Client struct {
	c *http.Client
}

// New returns a new HTTP client configured with the timeout from the application configuration.
// It sets up the client with the specified timeout and logs the timeout at debug level.
func New() *Client {
	timeout := time.Duration(config.Config.Timeout)
	slog.Debug(
		"create http client",
		"timeout", timeout,
	)
	return &Client{
		c: &http.Client{
			Timeout: timeout * time.Second,
		},
	}
}

// Do performs an HTTP request with logging and error handling.
// It logs request details and duration, handles non-2xx responses by reading the response body,
// and returns the response or an error.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := c.c.Do(req)
	if err != nil {
		slog.Error(
			"http request failed",
			"method", req.Method,
			"url", req.URL.String(),
			"error", err,
			"duration", time.Since(start),
		)
		return nil, err
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		err := resp.Body.Close()
		slog.Error(
			"http request returned error",
			"method", req.Method,
			"url", req.URL.String(),
			"status", resp.StatusCode,
			"body", string(body),
			"duration", time.Since(start),
		)
		if err != nil {
			return nil, fmt.Errorf("http response close failed: %w", err)
		}

		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}

	slog.Debug(
		"http request success",
		"method", req.Method,
		"url", req.URL,
		"status", resp.StatusCode,
		"duration", time.Since(start),
	)

	return resp, nil
}

// DoJSON performs an HTTP request and decodes the JSON response into the provided value.
// It calls Do to execute the request, then decodes the response body as JSON.
func (c *Client) DoJSON(req *http.Request, v any) error {
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("receive response failed: %w", err)
	}
	defer resp.Body.Close() // nolint
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("parse json failed: %w", err)
	}
	return nil
}
