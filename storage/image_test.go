package storage

import (
	"encoding/json"
	"testing"
)

func TestImageSizes_ToJson(t *testing.T) {
	imageSizes := ImageSizes{
		Original: Dimensions{Width: 500, Height: 300},
		Xs:       &Dimensions{Width: 300, Height: 100},
		S:        nil,
		M:        nil,
		L:        nil,
		XL:       nil,
		XXL:      nil,
		XXXL:     nil,
	}

	data, err := json.Marshal(imageSizes)
	if err != nil {
		t.Fatal("error converting to String", err)
	}

	expectedJson := "{\"original\":{\"width\":500,\"height\":300},\"xs\":{\"width\":300,\"height\":100}}"

	if string(data) != expectedJson {
		t.Fatalf("Expected: %s\nGot: %s", expectedJson, string(data))
	}
}
