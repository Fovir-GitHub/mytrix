package config

import (
	"testing"
)

func TestValidateGotify_Disabled(t *testing.T) {
	// Test when Gotify is disabled
	cfg := MytrixConfig{
		Gotify: GotifyConfig{
			Enabled: false,
			Server:  "", // Should not matter when disabled
			Token:   "", // Should not matter when disabled
		},
	}
	err := cfg.validateGotify()
	if err != nil {
		t.Fatalf("Expected no error when Gotify is disabled, got %v", err)
	}
}

func TestValidateGotify_Enabled_MissingServer(t *testing.T) {
	// Test when Gotify is enabled but server is missing
	cfg := MytrixConfig{
		Gotify: GotifyConfig{
			Enabled: true,
			Server:  "", // Missing server
			Token:   "valid-token",
		},
	}
	err := cfg.validateGotify()
	if err == nil {
		t.Fatalf("Expected error when Gotify is enabled but server is missing")
	}
	if err.Error() != "MYTRIX_GOTIFY_SERVER and MYTRIX_GOTIFY_TOKEN are required when MYTRIX_GOTIFY_ENABLE=true" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestValidateGotify_Enabled_MissingToken(t *testing.T) {
	// Test when Gotify is enabled but token is missing
	cfg := MytrixConfig{
		Gotify: GotifyConfig{
			Enabled: true,
			Server:  "valid-server",
			Token:   "", // Missing token
		},
	}
	err := cfg.validateGotify()
	if err == nil {
		t.Fatalf("Expected error when Gotify is enabled but token is missing")
	}
	if err.Error() != "MYTRIX_GOTIFY_SERVER and MYTRIX_GOTIFY_TOKEN are required when MYTRIX_GOTIFY_ENABLE=true" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestValidateGotify_Enabled_Valid(t *testing.T) {
	// Test when Gotify is enabled and both server and token are present
	cfg := MytrixConfig{
		Gotify: GotifyConfig{
			Enabled: true,
			Server:  "valid-server",
			Token:   "valid-token",
		},
	}
	err := cfg.validateGotify()
	if err != nil {
		t.Fatalf("Expected no error when Gotify is enabled with valid config, got %v", err)
	}
}

func TestValidate_NoErrors(t *testing.T) {
	// Test validate function with no errors
	cfg := MytrixConfig{
		Gotify: GotifyConfig{
			Enabled: false, // Disable to avoid validation errors
		},
		Wakapi: WakapiConfig{
			Enabled: false, // Disable to avoid validation errors
		},
		Umami: UmamiConfig{
			Enabled: false, // Disable to avoid validation errors
		},
	}
	err := cfg.validate()
	if err != nil {
		t.Fatalf("Expected no error with all services disabled, got %v", err)
	}
}

func TestValidate_WithErrors(t *testing.T) {
	// Test validate function with errors
	cfg := MytrixConfig{
		Gotify: GotifyConfig{
			Enabled: true,
			Server:  "", // Missing server
			Token:   "", // Missing token
		},
		Wakapi: WakapiConfig{
			Enabled: true,
			Server:  "", // Missing server
			APIKey:  "", // Missing API key
			UserID:  "", // Missing user ID
		},
		Umami: UmamiConfig{
			Enabled:  true,
			Server:   "", // Missing server
			Username: "", // Missing username
			Password: "", // Missing password
		},
	}
	err := cfg.validate()
	if err == nil {
		t.Fatalf("Expected error with missing required fields")
	}
	// Check that we got multiple errors joined together
	if err.Error() == "" {
		t.Fatalf("Expected non-empty error message")
	}
}
