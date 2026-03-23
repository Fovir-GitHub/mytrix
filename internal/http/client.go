package http

import (
	"net/http"
	"time"

	"github.com/Fovir-GitHub/mytrix/internal/config"
)

type Client struct {
	c *http.Client
}

func New() *Client {
	return &Client{
		c: &http.Client{
			Timeout: time.Duration(config.Config.Timeout) * time.Second,
		},
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.c.Do(req)
}
