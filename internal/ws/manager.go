package ws

import (
	"log/slog"
	"sync"

	"github.com/Fovir-GitHub/mytrix/internal/model"
)

type Manager struct {
	clients map[string]*Client
	events  chan *model.WsEvent
	once    sync.Once
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]*Client),
		events:  make(chan *model.WsEvent),
	}
}

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

func (m *Manager) AddIfEnabled(name, url string, enabled bool) {
	if enabled {
		m.Add(name, url)
	}
}

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
