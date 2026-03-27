package ws

import (
	"log/slog"
	"time"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/gorilla/websocket"
)

type Client struct {
	url  string
	conn *websocket.Conn
	recv chan []byte
}

func NewClient(url string) *Client {
	return &Client{
		url:  url,
		recv: make(chan []byte, config.Config.WS.RecvBufferSize),
	}
}

func (c *Client) Start() {
	go c.connectLoop()
}

func (c *Client) connectLoop() {
	retryInterval := time.Duration(config.Config.WS.RetryInterval)
	for {
		conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
		if err != nil {
			slog.Error(
				"websocket connect failed",
				"err", err,
			)
			time.Sleep(retryInterval * time.Second)
			continue
		}
		slog.Debug(
			"websocket connected",
			"url", c.url,
		)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				slog.Error(
					"websocket read error",
					"err", err,
				)
				conn.Close()
				break
			}
			slog.Debug(
				"websocket message received",
				"msg", string(msg),
			)
			c.recv <- msg
		}
		time.Sleep(retryInterval * time.Second)
	}
}

func (c *Client) Receive() <-chan []byte {
	return c.recv
}
