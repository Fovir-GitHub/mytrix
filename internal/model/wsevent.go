package model

// WsEvent represents a WebSocket event with its source and data payload.
type WsEvent struct {
	Source string
	Data   []byte
}
