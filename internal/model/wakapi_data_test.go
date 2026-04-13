package model

import (
	"testing"
)

func TestWakapiDataStructure(t *testing.T) {
	// Test that WakapiData can be created with correct values
	data := WakapiData{
		TotalTime:        "5 hrs 0 mins",
		ReadableInterval: "last 7 days",
		Langs: []WakapiLanguage{
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
		},
	}

	if data.TotalTime != "5 hrs 0 mins" {
		t.Errorf("Expected total time '5 hrs 0 mins', got %q", data.TotalTime)
	}
	if data.ReadableInterval != "last 7 days" {
		t.Errorf("Expected interval 'last 7 days', got %q", data.ReadableInterval)
	}
	if len(data.Langs) != 2 {
		t.Errorf("Expected 2 languages, got %d", len(data.Langs))
	}
}

func TestWakapiData_ToView(t *testing.T) {
	// Test the toView method
	data := WakapiData{
		TotalTime:        "5 hrs",
		ReadableInterval: "weekly",
		Langs: []WakapiLanguage{
			{
				Name:    "Go",
				Text:    "5 hrs",
				Percent: 100.0,
			},
		},
	}
	view := data.toView()

	if view == nil {
		t.Errorf("Expected non-nil view")
		return
	}
	if view.Interval != "weekly" {
		t.Errorf("Expected interval 'weekly', got %q", view.Interval)
	}
	if view.Total != "5 hrs" {
		t.Errorf("Expected total '5 hrs', got %q", view.Total)
	}
}
