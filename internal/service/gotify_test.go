package service

import (
	"testing"

	"github.com/Fovir-GitHub/mytrix/internal/model"
)

func TestGotifyService_HandleEvent(t *testing.T) {
	// Test successful unmarshaling
	service := &GotifyService{}
	eventData := []byte(`{"Title":"Test Title","Message":"Test Message","ID":123}`)
	event := &model.WsEvent{Data: eventData}

	msg, err := service.HandleEvent(event)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if msg == nil {
		t.Fatalf("Expected non-nil message")
	}
	if msg.Title != "Test Title" {
		t.Errorf("Expected Title 'Test Title', got '%s'", msg.Title)
	}
	if msg.Message != "Test Message" {
		t.Errorf("Expected Message 'Test Message', got '%s'", msg.Message)
	}
	if msg.ID != 123 {
		t.Errorf("Expected ID 123, got %d", msg.ID)
	}

	// Test invalid JSON
	invalidEvent := &model.WsEvent{Data: []byte(`invalid json`)}
	_, err = service.HandleEvent(invalidEvent)
	if err == nil {
		t.Fatalf("Expected error for invalid JSON")
	}
}
