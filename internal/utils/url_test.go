package utils_test

import (
	"testing"

	"codeberg.org/Fovir/mytrix/internal/utils"
)

func TestParseURLHost(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		raw  string
		want string
	}{
		{
			name: "Host Only",
			raw:  "https://example.com",
			want: "example.com",
		},
		{
			name: "Host with Path",
			raw:  "https://example.com/path/to/file",
			want: "example.com",
		},
		{
			name: "Query",
			raw:  "http://example.com/product?id=123",
			want: "example.com",
		},
		{
			name: "2-level Domain",
			raw:  "https://www.example.com/path/to/file?id=123",
			want: "www.example.com",
		},
		{
			name: "N-level Domain",
			raw:  "ws://1.2.3.4.5.6.7.8.9.10.example.com/path/to/file?id=123",
			want: "1.2.3.4.5.6.7.8.9.10.example.com",
		},
		{
			name: "Network Name",
			raw:  "http://example",
			want: "example",
		},
		{
			name: "Invalid Host",
			raw:  "example",
			want: "",
		},
		{
			name: "Empty URL",
			raw:  "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.ParseURLHost(tt.raw)

			if got != tt.want {
				t.Errorf("ParseURLHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
