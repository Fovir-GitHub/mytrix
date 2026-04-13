package service

import (
	"errors"
	"testing"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/http"
	"codeberg.org/Fovir/mytrix/internal/model"
	"codeberg.org/Fovir/mytrix/internal/scheduler"
)

func TestNewWakapiService_Disabled(t *testing.T) {
	// Initialize config if needed
	if config.Config == nil {
		config.Config = &config.MytrixConfig{}
	}

	// Temporarily disable Wakapi
	originalEnabled := config.Config.Wakapi.Enabled
	config.Config.Wakapi.Enabled = false
	defer func() { config.Config.Wakapi.Enabled = originalEnabled }()

	service := newWakapiService(nil, nil)
	if service == nil {
		t.Fatalf("Expected non-nil service")
	}
	// Check that it's a NoopWakapiService
	if _, ok := service.(*NoopWakapiService); !ok {
		t.Fatalf("Expected NoopWakapiService when Wakapi is disabled")
	}
}

func TestNewWakapiService_Enabled(t *testing.T) {
	// Initialize config if needed
	if config.Config == nil {
		config.Config = &config.MytrixConfig{}
	}

	// Temporarily enable Wakapi
	originalEnabled := config.Config.Wakapi.Enabled
	originalServer := config.Config.Wakapi.Server
	originalAPIKey := config.Config.Wakapi.APIKey
	originalUserID := config.Config.Wakapi.UserID
	config.Config.Wakapi.Enabled = true
	config.Config.Wakapi.Server = "test.server"
	config.Config.Wakapi.APIKey = "test-key"
	config.Config.Wakapi.UserID = "test-user"
	defer func() {
		config.Config.Wakapi.Enabled = originalEnabled
		config.Config.Wakapi.Server = originalServer
		config.Config.Wakapi.APIKey = originalAPIKey
		config.Config.Wakapi.UserID = originalUserID
	}()

	service := newWakapiService(&http.Client{}, &scheduler.Scheduler{})
	if service == nil {
		t.Fatalf("Expected non-nil service")
	}
	// Check that it's a RealWakapiService
	if _, ok := service.(*RealWakapiService); !ok {
		t.Fatalf("Expected RealWakapiService when Wakapi is enabled")
	}
}

func TestNoopWakapiService_FetchData(t *testing.T) {
	testErr := errors.New("test error")
	service := &NoopWakapiService{err: testErr}
	data, err := service.FetchReport(model.WakapiIntervalToday)
	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
	if data != "" {
		t.Errorf("Expected empty string, got %q", data)
	}
}

func TestNoopWakapiService_FetchReport(t *testing.T) {
	testErr := errors.New("test error")
	service := &NoopWakapiService{err: testErr}
	_, err := service.fetchData(model.WakapiIntervalToday)
	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
}
