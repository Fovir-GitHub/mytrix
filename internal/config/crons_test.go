package config

import (
	"testing"
)

func TestValidateCrons_ValidCrons(t *testing.T) {
	// Test with valid cron expressions
	cfg := MytrixConfig{}
	crons := []string{
		"0 9 * * *", // Every day at 9:00
		"0 9 * * 1", // Every Monday at 9:00
		"0 9 1 * *", // First day of month at 9:00
		"0 9 1 1 *", // January 1st at 9:00
	}
	err := cfg.validateCrons(crons)
	if err != nil {
		t.Fatalf("Expected no error for valid crons, got %v", err)
	}
}

func TestValidateCrons_InvalidCrons(t *testing.T) {
	// Test with invalid cron expression
	cfg := MytrixConfig{}
	crons := []string{"invalid cron"}
	err := cfg.validateCrons(crons)
	if err == nil {
		t.Fatalf("Expected error for invalid cron expression")
	}
}

func TestValidateCrons_MixedValidInvalid(t *testing.T) {
	// Test with mix of valid and invalid crons
	cfg := MytrixConfig{}
	crons := []string{
		"0 9 * * *", // Valid
		"invalid",   // Invalid
		"0 9 1 * *", // Valid
	}
	err := cfg.validateCrons(crons)
	if err == nil {
		t.Fatalf("Expected error when mix of valid and invalid crons")
	}
}

func TestValidateCrons_Empty(t *testing.T) {
	// Test with empty cron list
	cfg := MytrixConfig{}
	crons := []string{}
	err := cfg.validateCrons(crons)
	if err != nil {
		t.Fatalf("Expected no error for empty cron list, got %v", err)
	}
}

func TestValidateCrons_SpecialExpressions(t *testing.T) {
	// Test with special cron expressions
	cfg := MytrixConfig{}
	crons := []string{
		"@yearly",
		"@monthly",
		"@weekly",
		"@daily",
		"@hourly",
	}
	err := cfg.validateCrons(crons)
	if err != nil {
		t.Fatalf("Expected no error for special expressions, got %v", err)
	}
}
