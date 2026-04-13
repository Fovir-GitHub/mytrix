package model

import (
	"strings"
	"testing"
)

// Note: templates are initialized in umami_test.go init()

func TestWakapiLanguageStructure(t *testing.T) {
	// Test that WakapiLanguage can be created with correct values
	lang := WakapiLanguage{
		Name:    "Go",
		Text:    "2 hrs 30 mins",
		Percent: 45.5,
	}

	if lang.Name != "Go" {
		t.Errorf("Expected name 'Go', got %q", lang.Name)
	}
	if lang.Text != "2 hrs 30 mins" {
		t.Errorf("Expected text '2 hrs 30 mins', got %q", lang.Text)
	}
	if lang.Percent != 45.5 {
		t.Errorf("Expected percent 45.5, got %f", lang.Percent)
	}
}

func TestGenerateLangReport_Empty(t *testing.T) {
	// Test with empty language list
	langs := []WakapiLanguage{}
	result := generateLangReport(langs)

	if result != "" {
		t.Errorf("Expected empty string for empty language list, got %q", result)
	}
}

func TestGenerateLangReport_SingleLanguage(t *testing.T) {
	// Test with single language - should contain code block markers
	langs := []WakapiLanguage{
		{
			Name:    "Go",
			Text:    "2 hrs 30 mins",
			Percent: 100.0,
		},
	}
	result := generateLangReport(langs)

	if !strings.Contains(result, "```text") {
		t.Errorf("Expected code block markers in output, got %q", result)
	}
	if !strings.Contains(result, "```") {
		t.Errorf("Expected closing code block marker in output, got %q", result)
	}
}

func TestGenerateLangReport_MultipleLanguages(t *testing.T) {
	// Test with multiple languages
	langs := []WakapiLanguage{
		{
			Name:    "Go",
			Text:    "2 hrs 30 mins",
			Percent: 50.0,
		},
		{
			Name:    "Python",
			Text:    "2 hrs 30 mins",
			Percent: 50.0,
		},
	}
	result := generateLangReport(langs)

	if !strings.Contains(result, "```text") {
		t.Errorf("Expected code block markers in output")
	}
	// Result should have content for both languages
	if len(result) == 0 {
		t.Errorf("Expected non-empty result for multiple languages")
	}
}
