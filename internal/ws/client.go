package ws

import (
	"log/slog"
	"sync"
	"time"

	"codeberg.org/Fovir/mytrix/internal/config"
	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection that automatically reconnects on failure.
type Client struct {
	url  string
	conn *websocket.Conn
	recv chan []byte
	done chan struct{}
	mu   sync.Mutex
}

// NewClient creates a new WebSocket client for the given URL.
func NewClient(url string) *Client {
	return &Client{
		url:  url,
		recv: make(chan []byte, config.Config.WS.RecvBufferSize),
		done: make(chan struct{}),
	}
}

// Start begins the connection and message reading loop in a background goroutine.
func (c *Client) Start() {
	go c.connectLoop()
}

// Stop closes the client connection and stops the background loops.
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
		slog.Info(
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
			slog.Error("websocket read error",
				"err", err)
			if err := conn.Close(); err != nil {
				slog.Error("websocket connection close failed", "err", err)
			}
			return true
		}
		slog.Info("websocket message received",
			"len", len(msg))
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

// Receive returns a receive-only channel for WebSocket messages.
func (c *Client) Receive() <-chan []byte {
	return c.recv
}
