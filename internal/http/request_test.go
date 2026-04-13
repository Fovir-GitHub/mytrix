package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/Fovir/mytrix/internal/config"
)

func init() {
	// Initialize config for tests
	if config.Config == nil {
		config.Config = &config.MytrixConfig{
			Timeout: 10,
		}
	}
}

func TestNewRequest_GET(t *testing.T) {
	// Test creating a GET request
	req, err := NewRequest(MethodGet, "https://example.com/api", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if req == nil {
		t.Fatalf("Expected non-nil request")
	}
	if req.Method != "GET" {
		t.Errorf("Expected method GET, got %q", req.Method)
	}
	if req.URL.String() != "https://example.com/api" {
		t.Errorf("Expected URL https://example.com/api, got %q", req.URL.String())
	}
}

func TestNewRequest_POST(t *testing.T) {
	// Test creating a POST request with body
	body := bytes.NewBufferString(`{"key": "value"}`)
	req, err := NewRequest(MethodPost, "https://example.com/api", body, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if req.Method != "POST" {
		t.Errorf("Expected method POST, got %q", req.Method)
	}
}

func TestNewRequest_WithHeaders(t *testing.T) {
	// Test creating request with headers
	headers := map[string]string{
		"Authorization": "Bearer token123",
		"Content-Type":  "application/json",
	}
	req, err := NewRequest(MethodGet, "https://example.com/api", nil, headers)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if req.Header.Get("Authorization") != "Bearer token123" {
		t.Errorf("Expected Authorization header not set correctly")
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header not set correctly")
	}
}

func TestNewRequest_InvalidURL(t *testing.T) {
	// Test with invalid URL
	req, err := NewRequest(MethodGet, "invalid://url with spaces", nil, nil)

	if err == nil {
		t.Fatalf("Expected error for invalid URL")
	}
	if req != nil {
		t.Fatalf("Expected nil request for invalid URL")
	}
}

func TestClient_New(t *testing.T) {
	// Test creating a new client
	// Save original timeout
	originalTimeout := config.Config.Timeout
	defer func() { config.Config.Timeout = originalTimeout }()

	config.Config.Timeout = 30
	client := New()

	if client == nil {
		t.Fatalf("Expected non-nil client")
	}
	if client.c == nil {
		t.Fatalf("Expected non-nil underlying http client")
	}
}

func TestClient_Do_Success(t *testing.T) {
	// Test successful request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`)) //nolint
	}))
	defer server.Close()

	// Save original timeout
	originalTimeout := config.Config.Timeout
	defer func() { config.Config.Timeout = originalTimeout }()

	config.Config.Timeout = 10
	client := New()

	req, err := NewRequest(MethodGet, server.URL, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestClient_Do_ErrorStatus(t *testing.T) {
	// Test request with error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`)) //nolint
	}))
	defer server.Close()

	// Save original timeout
	originalTimeout := config.Config.Timeout
	defer func() { config.Config.Timeout = originalTimeout }()

	config.Config.Timeout = 10
	client := New()

	req, err := NewRequest(MethodGet, server.URL, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	_, err = client.Do(req)
	if err == nil {
		t.Fatalf("Expected error for 404 status")
	}
}

func TestClient_DoJSON(t *testing.T) {
	// Test JSON decoding
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "test", "value": 42}`)) //nolint
	}))
	defer server.Close()

	// Save original timeout
	originalTimeout := config.Config.Timeout
	defer func() { config.Config.Timeout = originalTimeout }()

	config.Config.Timeout = 10
	client := New()

	req, err := NewRequest(MethodGet, server.URL, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]any
	err = client.DoJSON(req, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result["name"] != "test" {
		t.Errorf("Expected name 'test', got %v", result["name"])
	}
	if result["value"] != float64(42) {
		t.Errorf("Expected value 42, got %v", result["value"])
	}
}

func TestClient_DoJSON_InvalidJSON(t *testing.T) {
	// Test JSON decoding with invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`)) //nolint
	}))
	defer server.Close()

	// Save original timeout
	originalTimeout := config.Config.Timeout
	defer func() { config.Config.Timeout = originalTimeout }()

	config.Config.Timeout = 10
	client := New()

	req, err := NewRequest(MethodGet, server.URL, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]interface{}
	err = client.DoJSON(req, &result)
	if err == nil {
		t.Fatalf("Expected error for invalid JSON")
	}
}
