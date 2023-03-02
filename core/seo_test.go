package core

import (
	"testing"
)

func TestFormatForSeo(t *testing.T) {
	values := []struct {
		Name     string
		Value    string
		Expected string
	}{
		{
			Name:     "Badly formatted data",
			Value:    "_some1 ran----dom __test -- to add-5- ",
			Expected: "some1-ran-dom-test-to-add-5",
		},
		{
			Name:     "Only bad data",
			Value:    " -_--__- _ --_ ",
			Expected: "",
		},
		{
			Name:     "Empty string",
			Value:    "",
			Expected: "",
		},
		{
			Name:     "Standard text to expect",
			Value:    "world war 2 plane ",
			Expected: "world-war-2-plane",
		},
	}

	for _, data := range values {
		t.Run(data.Name, func(t *testing.T) {
			result := FormatForSeo(data.Value)
			if result != data.Expected {
				t.Fatalf("Expected %s, got %s\n", data.Expected, result)
			}
		})
	}
}
