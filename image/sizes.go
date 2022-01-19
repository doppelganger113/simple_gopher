package image

import "fmt"

type Dimensions struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
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

type Sizes struct {
	Original Dimensions  `json:"original"`
	Xs       *Dimensions `json:"xs,omitempty"`
	S        *Dimensions `json:"s,omitempty"`
	M        *Dimensions `json:"m,omitempty"`
	L        *Dimensions `json:"l,omitempty"`
	XL       *Dimensions `json:"xl,omitempty"`
	XXL      *Dimensions `json:"xxl,omitempty"`
	XXXL     *Dimensions `json:"xxxl,omitempty"`
}

func (sizes Sizes) ToString() string {
	return fmt.Sprintf(
		"Sizes{original: %s, xs: %s, s: %s, m: %s, l: %s, xl: %s, xxl: %s, xxxl: %s}",
		sizes.Original.ToString(),
		sizes.Xs.ToString(),
		sizes.S.ToString(),
		sizes.M.ToString(),
		sizes.L.ToString(),
		sizes.XL.ToString(),
		sizes.XXL.ToString(),
		sizes.XXXL.ToString(),
	)
}

func (sizes Sizes) IsEqualTo(compareSizes Sizes) bool {
	if !AreEqual(&sizes.Original, &compareSizes.Original) {
		return false
	}
	if !AreEqual(sizes.Xs, compareSizes.Xs) {
		return false
	}
	if !AreEqual(sizes.S, compareSizes.S) {
		return false
	}
	if !AreEqual(sizes.M, compareSizes.M) {
		return false
	}
	if !AreEqual(sizes.L, compareSizes.L) {
		return false
	}
	if !AreEqual(sizes.XL, compareSizes.XL) {
		return false
	}
	if !AreEqual(sizes.XXL, compareSizes.XXL) {
		return false
	}
	if !AreEqual(sizes.XXL, compareSizes.XXL) {
		return false
	}
	if !AreEqual(sizes.XXXL, compareSizes.XXXL) {
		return false
	}
	return true
}

func (sizes Sizes) GetAllDimensions() []Dimensions {
	dimensions := []Dimensions{sizes.Original}

	if sizes.Xs != nil {
		dimensions = append(dimensions, *sizes.Xs)
	}
	if sizes.S != nil {
		dimensions = append(dimensions, *sizes.S)
	}
	if sizes.M != nil {
		dimensions = append(dimensions, *sizes.M)
	}
	if sizes.L != nil {
		dimensions = append(dimensions, *sizes.L)
	}
	if sizes.XL != nil {
		dimensions = append(dimensions, *sizes.XL)
	}
	if sizes.XXL != nil {
		dimensions = append(dimensions, *sizes.XXL)
	}
	if sizes.XXXL != nil {
		dimensions = append(dimensions, *sizes.XXXL)
	}

	return dimensions
}
