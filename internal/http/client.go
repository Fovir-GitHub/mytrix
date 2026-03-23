package http

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Fovir-GitHub/mytrix/internal/config"
)

type Client struct {
	c *http.Client
}

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
		resp.Body.Close()

		slog.Error(
			"http request returned error",
			"method", req.Method,
			"url", req.URL.String(),
			"status", resp.StatusCode,
			"body", string(body),
			"duration", time.Since(start),
		)

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
