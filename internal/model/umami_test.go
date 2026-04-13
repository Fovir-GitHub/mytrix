package model

import (
	"testing"

	"github.com/Fovir-GitHub/mytrix/internal/config"
)

func init() {
	// Initialize config and templates for tests
	config.Config = &config.MytrixConfig{}
	InitTemplates()
}

func TestUmamiWebsiteStructure(t *testing.T) {
	// Test that UmamiWebsite can be created with correct values
	website := &UmamiWebsite{
		ID:     "1",
		Name:   "My Site",
		Domain: "example.com",
		Stat: &UmamiWebsiteStat{
			Visitors: 100,
			Visits:   150,
			Bounces:  30,
		},
	}

	if website.ID != "1" {
		t.Errorf("Expected ID '1', got %q", website.ID)
	}
	if website.Name != "My Site" {
		t.Errorf("Expected name 'My Site', got %q", website.Name)
	}
	if website.Domain != "example.com" {
		t.Errorf("Expected domain 'example.com', got %q", website.Domain)
	}
	if website.Stat.Visitors != 100 {
		t.Errorf("Expected 100 visitors, got %d", website.Stat.Visitors)
	}
}

func TestUmamiWebsite_ToView(t *testing.T) {
	// Test the toView method for bounce rate calculation
	website := &UmamiWebsite{
		ID:     "1",
		Name:   "Site",
		Domain: "example.com",
		Stat: &UmamiWebsiteStat{
			Visitors: 100,
			Visits:   200,
			Bounces:  50,
		},
	}
	view := website.toView()

	if view == nil {
		t.Errorf("Expected non-nil view")
		return
	}
	if view.Name != "Site" {
		t.Errorf("Expected name 'Site', got %q", view.Name)
	}
	if view.Domain != "example.com" {
		t.Errorf("Expected domain 'example.com', got %q", view.Domain)
	}
	if view.Visitors != 100 {
		t.Errorf("Expected 100 visitors, got %d", view.Visitors)
	}
	if view.Visits != 200 {
		t.Errorf("Expected 200 visits, got %d", view.Visits)
	}
	// Bounce rate: 50/200 * 100 = 25.00%
	if view.BouncesRate != "25.00%" {
		t.Errorf("Expected bounce rate '25.00%%', got %q", view.BouncesRate)
	}
}

func TestUmamiWebsite_ToView_ZeroVisits(t *testing.T) {
	// Test toView with zero visits
	website := &UmamiWebsite{
		ID:     "1",
		Name:   "Site",
		Domain: "test.com",
		Stat: &UmamiWebsiteStat{
			Visitors: 0,
			Visits:   0,
			Bounces:  0,
		},
	}
	view := website.toView()

	if view == nil {
		t.Errorf("Expected non-nil view")
		return
	}
	if view.BouncesRate != "0%" {
		t.Errorf("Expected bounce rate '0%%' when zero visits, got %q", view.BouncesRate)
	}
}
