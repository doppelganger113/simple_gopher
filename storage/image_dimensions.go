package storage

import "fmt"

type Dimensions struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

func AreEqual(dimensions *Dimensions, d *Dimensions) bool {
	if dimensions == nil && d == nil {
		return true
	}

	if dimensions != nil && d != nil {
		if dimensions.Height != d.Height {
			return false
		}
		if dimensions.Width != d.Width {
			return false
		}
		return true
	}

	return false
}

func (dimensions *Dimensions) ToString() string {
	if dimensions == nil {
		return "nil"
	}
	return fmt.Sprintf(
		"Dimensions{width: %d, height: %d}",
		dimensions.Width,
		dimensions.Height,
	)
}
