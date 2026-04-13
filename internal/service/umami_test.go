package service

import (
	"errors"
	"testing"

	"codeberg.org/Fovir/mytrix/internal/config"
)

func TestNewUmamiService_Disabled(t *testing.T) {
	// Initialize config if needed
	if config.Config == nil {
		config.Config = &config.MytrixConfig{}
	}

	// Temporarily disable Umami
	originalEnabled := config.Config.Umami.Enabled
	config.Config.Umami.Enabled = false
	defer func() { config.Config.Umami.Enabled = originalEnabled }()

	service := newUmamiService(nil)
	if service == nil {
		t.Fatalf("Expected non-nil service")
	}
	// Check that it's a NoopUmamiService
	if _, ok := service.(*NoopUmamiService); !ok {
		t.Fatalf("Expected NoopUmamiService when Umami is disabled")
	}
}

func TestNoopUmamiService_FetchReport(t *testing.T) {
	testErr := errors.New("test error")
	service := &NoopUmamiService{err: testErr}
	result := service.FetchReport(nil)
	if result != testErr.Error() {
		t.Errorf("Expected error %v, got %v", testErr.Error(), result)
	}
}
