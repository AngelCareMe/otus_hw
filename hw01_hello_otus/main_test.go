package main

import (
	"golang.org/x/example/hello/reverse"
	"testing"
)

func TestReverseString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, OTUS!", "!SUTO ,olleH"},
		{"", ""},
		{"abc", "cba"},
		{"12345", "54321"},
		{"  whitespace  ", "  ecapsetihw  "},
	}

	for _, test := range tests {
		if result := reverse.String(test.input); result != test.expected {
			t.Errorf("reverse.String(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
