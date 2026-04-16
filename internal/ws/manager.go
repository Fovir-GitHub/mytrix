package ws

import (
	"log/slog"
	"sync"

	"codeberg.org/Fovir/mytrix/internal/model"
)

// Manager manages multiple WebSocket client connections and aggregates their events into a single channel.
type Manager struct {
	clients map[string]*Client
	events  chan *model.WsEvent
	once    sync.Once
}

// NewManager creates a new Manager with an empty client map and event channel.
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]*Client),
		events:  make(chan *model.WsEvent),
	}
}

// Add creates a new WebSocket client, registers it by name, and starts it.
func (m *Manager) Add(name, url string) {
	client := NewClient(url)
	m.clients[name] = client
	client.Start()
	slog.Info(
		"websocket added",
		"name", name,
		"url", url,
	)
}

// AddIfEnabled conditionally adds a WebSocket client if the enabled flag is true.
func (m *Manager) AddIfEnabled(name, url string, enabled bool) {
	if enabled {
		m.Add(name, url)
	}
}

// Events returns a receive-only channel of WebSocket events from all managed clients.
// All client receive channels are multiplexed into this single channel.
func (m *Manager) Events() <-chan *model.WsEvent {
	m.once.Do(func() {
		for name, client := range m.clients {
			go func(name string, c *Client) {
				for msg := range c.Receive() {
					select {
					case m.events <- &model.WsEvent{
						Source: name,
						Data:   msg,
					}:
					default:
					}
				}
			}(name, client)
		}
	})
	return m.events
}
