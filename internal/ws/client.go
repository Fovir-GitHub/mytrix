package ws

import (
	"log/slog"
	"sync"
	"time"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/gorilla/websocket"
)

type Client struct {
	url  string
	conn *websocket.Conn
	recv chan []byte
	done chan struct{}
	mu   sync.Mutex
}

func NewClient(url string) *Client {
	return &Client{
		url:  url,
		recv: make(chan []byte, config.Config.WS.RecvBufferSize),
		done: make(chan struct{}),
	}
}

func (c *Client) Start() {
	go c.connectLoop()
}

func (c *Client) Stop() {
	close(c.done)
	c.mu.Lock()
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			slog.Error("websocket stop close failed", "err", err)
		}
	}
	c.mu.Unlock()
}

func (c *Client) connectLoop() {
	retryInterval := time.Duration(config.Config.WS.RetryInterval)
	for {
		select {
		case <-c.done:
			return
		default:
		}

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

		c.mu.Lock()
		c.conn = conn
		c.mu.Unlock()

		if !c.readLoop(conn) {
			return
		}
		time.Sleep(retryInterval * time.Second)
	}
}

func (c *Client) readLoop(conn *websocket.Conn) bool {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Error(
				"websocket read error",
				"err", err,
			)
			if err := conn.Close(); err != nil {
				slog.Error("websocket connection close failed", "err", err)
			}
			return true
		}
		slog.Debug(
			"websocket message received",
			"msg", string(msg),
		)
		select {
		case c.recv <- msg:
		case <-c.done:
			if err := conn.Close(); err != nil {
				slog.Error("websocket connection close failed", "err", err)
			}
			return false
		default:
			slog.Warn("recv channel full, dropping message")
		}
	}
}

func (c *Client) Receive() <-chan []byte {
	return c.recv
}
