package http_util

import (
	"testing"
)

func TestGetTokenFromHeader(t *testing.T) {
	result := GetTokenFromHeader("Bearer token")
	if result != "token" {
		t.Errorf("expected %s to equal token", result)
	}

	result = GetTokenFromHeader("Token")
	if result != "" {
		t.Errorf("expected result to be empty, got %s", result)
	}
}

func TestToUint(t *testing.T) {
	data := []struct {
		value    string
		expected uint
	}{
		{value: "", expected: 0},
		{value: "0", expected: 0},
		{value: "2", expected: 2},
		{value: "-2", expected: 0},
	}

	for _, d := range data {
		result := ToUint(d.value)
		if result != d.expected {
			t.Errorf("Expected '%s' to be %d, got %d", d.value, d.expected, result)
		}
	}
}
