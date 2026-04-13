package utils

import (
	"testing"
)

func TestFilter_WithIntegers(t *testing.T) {
	// Test filtering integers greater than 5
	numbers := []int{1, 5, 10, 3, 8, 2}
	result := Filter(numbers, func(n *int) bool {
		return *n > 5
	})

	if len(result) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(result))
	}
	if result[0] != 10 || result[1] != 8 {
		t.Errorf("Expected [10, 8], got %v", result)
	}
}

func TestFilter_WithStrings(t *testing.T) {
	// Test filtering strings with length > 2
	words := []string{"a", "ab", "abc", "abcd", "x"}
	result := Filter(words, func(s *string) bool {
		return len(*s) > 2
	})

	if len(result) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(result))
	}
	if result[0] != "abc" || result[1] != "abcd" {
		t.Errorf("Expected [abc, abcd], got %v", result)
	}
}

func TestFilter_EmptySlice(t *testing.T) {
	// Test filtering empty slice
	numbers := []int{}
	result := Filter(numbers, func(n *int) bool {
		return *n > 5
	})

	if len(result) != 0 {
		t.Errorf("Expected 0 elements, got %d", len(result))
	}
}

func TestFilter_NoMatches(t *testing.T) {
	// Test when no elements match the predicate
	numbers := []int{1, 2, 3, 4, 5}
	result := Filter(numbers, func(n *int) bool {
		return *n > 10
	})

	if len(result) != 0 {
		t.Errorf("Expected 0 elements, got %d", len(result))
	}
}

func TestFilter_AllMatches(t *testing.T) {
	// Test when all elements match the predicate
	numbers := []int{1, 2, 3, 4, 5}
	result := Filter(numbers, func(n *int) bool {
		return *n > 0
	})

	if len(result) != 5 {
		t.Errorf("Expected 5 elements, got %d", len(result))
	}
}
