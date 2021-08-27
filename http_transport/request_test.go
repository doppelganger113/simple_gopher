package http_transport

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTokenFromHeader(t *testing.T) {
	result := GetTokenFromHeader("Bearer token")
	assert.Equal(t, result, "token")

	result = GetTokenFromHeader("Token")
	assert.Empty(t, result)
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
		errMsg := fmt.Sprintf("Expected '%s' to be %d, got %d", d.value, d.expected, result)
		assert.Equal(t, result, d.expected, errMsg)
	}
}
