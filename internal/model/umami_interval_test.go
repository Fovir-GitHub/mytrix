package model

import (
	"testing"
	"time"
)

func TestGetStartDate(t *testing.T) {
	// Test that start date is set to midnight
	testTime := time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)
	result := getStartDate(&testTime)

	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
		t.Errorf("Expected start of day (00:00:00), got %02d:%02d:%02d", result.Hour(), result.Minute(), result.Second())
	}
	if result.Year() != 2023 || result.Month() != 5 || result.Day() != 15 {
		t.Errorf("Expected date 2023-05-15, got %d-%02d-%02d", result.Year(), result.Month(), result.Day())
	}
}

func TestGetEndDate(t *testing.T) {
	// Test that end date is set to end of day
	testTime := time.Date(2023, 5, 15, 10, 20, 30, 0, time.UTC)
	result := getEndDate(&testTime)

	if result.Hour() != 23 || result.Minute() != 59 || result.Second() != 59 {
		t.Errorf("Expected end of day (23:59:59), got %02d:%02d:%02d", result.Hour(), result.Minute(), result.Second())
	}
	if result.Year() != 2023 || result.Month() != 5 || result.Day() != 15 {
		t.Errorf("Expected date 2023-05-15, got %d-%02d-%02d", result.Year(), result.Month(), result.Day())
	}
}

func TestUmamiIntervalYesterday(t *testing.T) {
	// Test yesterday interval
	interval := UmamiIntervalYesterday()

	if interval == nil {
		t.Errorf("Expected non-nil interval")
		return
	}
	if interval.Start.After(interval.End) {
		t.Errorf("Start date should not be after end date")
	}

	// Start should be at midnight, end should be at 23:59:59
	if interval.Start.Hour() != 0 || interval.Start.Minute() != 0 {
		t.Errorf("Start should be at midnight")
	}
	if interval.End.Hour() != 23 || interval.End.Minute() != 59 {
		t.Errorf("End should be at 23:59:59")
	}
}

func TestUmamiIntervalLastWeek(t *testing.T) {
	// Test last week interval
	interval := UmamiIntervalLastWeek()

	if interval == nil {
		t.Errorf("Expected non-nil interval")
		return
	}
	if interval.Start.After(interval.End) {
		t.Errorf("Start date should not be after end date")
	}

	// Duration should be approximately 7 days
	duration := interval.End.Sub(interval.Start)
	days := duration.Hours() / 24
	if days < 6 || days > 8 {
		t.Errorf("Expected week interval to be ~7 days, got %f days", days)
	}
}

func TestUmamiIntervalLastMonth(t *testing.T) {
	// Test last month interval
	interval := UmamiIntervalLastMonth()

	if interval == nil {
		t.Errorf("Expected non-nil interval")
		return
	}
	if interval.Start.After(interval.End) {
		t.Errorf("Start date should not be after end date")
	}

	// Duration should be approximately 30 days
	duration := interval.End.Sub(interval.Start)
	days := duration.Hours() / 24
	if days < 25 || days > 35 {
		t.Errorf("Expected month interval to be ~30 days, got %f days", days)
	}
}

func TestUmamiIntervalLastYear(t *testing.T) {
	// Test last year interval
	interval := UmamiIntervalLastYear()

	if interval == nil {
		t.Errorf("Expected non-nil interval")
		return
	}
	if interval.Start.After(interval.End) {
		t.Errorf("Start date should not be after end date")
	}

	// Duration should be approximately 365 days
	duration := interval.End.Sub(interval.Start)
	days := duration.Hours() / 24
	if days < 360 || days > 370 {
		t.Errorf("Expected year interval to be ~365 days, got %f days", days)
	}
}

func TestParseUmamiInterval_Daily(t *testing.T) {
	// Test parsing "daily"
	interval, err := ParseUmamiInterval("daily")
	if err != nil {
		t.Errorf("Expected no error for 'daily', got %v", err)
	}
	if interval == nil {
		t.Errorf("Expected non-nil interval for 'daily'")
	}
}

func TestParseUmamiInterval_CaseInsensitive(t *testing.T) {
	// Test case insensitivity
	tests := []string{"DAILY", "Daily", "dAiLy"}
	for _, test := range tests {
		interval, err := ParseUmamiInterval(test)
		if err != nil {
			t.Errorf("Expected no error for %q, got %v", test, err)
		}
		if interval == nil {
			t.Errorf("Expected non-nil interval for %q", test)
		}
	}
}

func TestParseUmamiInterval_WithWhitespace(t *testing.T) {
	// Test parsing with whitespace
	interval, err := ParseUmamiInterval("  daily  ")
	if err != nil {
		t.Errorf("Expected no error for whitespace-padded 'daily', got %v", err)
	}
	if interval == nil {
		t.Errorf("Expected non-nil interval for whitespace-padded 'daily'")
	}
}

func TestParseUmamiInterval_Invalid(t *testing.T) {
	// Test parsing invalid interval
	interval, err := ParseUmamiInterval("invalid")

	if err == nil {
		t.Errorf("Expected error for invalid interval")
	}
	if interval != nil {
		t.Errorf("Expected nil interval for invalid input")
	}
}

func TestParseUmamiInterval_AllValidIntervals(t *testing.T) {
	// Test all valid intervals
	validIntervals := []string{"daily", "weekly", "monthly", "yearly"}
	for _, intervalStr := range validIntervals {
		interval, err := ParseUmamiInterval(intervalStr)
		if err != nil {
			t.Errorf("Expected no error for %q, got %v", intervalStr, err)
		}
		if interval == nil {
			t.Errorf("Expected non-nil interval for %q", intervalStr)
		}
	}
}
