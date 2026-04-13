package feed

import (
	"testing"
)

func TestNew(t *testing.T) {
	// Test parser creation
	parser := New()

	if parser == nil {
		t.Fatalf("Expected non-nil parser")
	}
	if parser.p == nil {
		t.Fatalf("Expected non-nil internal gofeed parser")
	}
}

func TestParser_ParseURL_InvalidURL(t *testing.T) {
	// Test parsing invalid URL
	parser := New()
	feed, items, err := parser.ParseURL("invalid://url")

	if err == nil {
		t.Fatalf("Expected error for invalid URL")
	}
	if feed != nil {
		t.Fatalf("Expected nil feed for invalid URL")
	}
	if items != nil {
		t.Fatalf("Expected nil items for invalid URL")
	}
}

func TestParser_ParseURL_Unreachable(t *testing.T) {
	// Test parsing unreachable URL
	parser := New()
	feed, items, err := parser.ParseURL("http://unreachable.invalid/feed.xml")

	if err == nil {
		t.Fatalf("Expected error for unreachable URL")
	}
	if feed != nil {
		t.Fatalf("Expected nil feed for unreachable URL")
	}
	if items != nil {
		t.Fatalf("Expected nil items for unreachable URL")
	}
}
