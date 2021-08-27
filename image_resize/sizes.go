package image_resize

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

type ImageSizes struct {
	Original Dimensions  `json:"original"`
	Xs       *Dimensions `json:"xs,omitempty"`
	S        *Dimensions `json:"s,omitempty"`
	M        *Dimensions `json:"m,omitempty"`
	L        *Dimensions `json:"l,omitempty"`
	XL       *Dimensions `json:"xl,omitempty"`
	XXL      *Dimensions `json:"xxl,omitempty"`
	XXXL     *Dimensions `json:"xxxl,omitempty"`
}

func (imageSizes ImageSizes) ToString() string {
	return fmt.Sprintf(
		"ImageSizes{original: %s, xs: %s, s: %s, m: %s, l: %s, xl: %s, xxl: %s, xxxl: %s}",
		imageSizes.Original.ToString(),
		imageSizes.Xs.ToString(),
		imageSizes.S.ToString(),
		imageSizes.M.ToString(),
		imageSizes.L.ToString(),
		imageSizes.XL.ToString(),
		imageSizes.XXL.ToString(),
		imageSizes.XXXL.ToString(),
	)
}

func (imageSizes ImageSizes) IsEqualTo(sizes ImageSizes) bool {
	if AreEqual(&imageSizes.Original, &sizes.Original) == false {
		return false
	}
	if AreEqual(imageSizes.Xs, sizes.Xs) == false {
		return false
	}
	if AreEqual(imageSizes.S, sizes.S) == false {
		return false
	}
	if AreEqual(imageSizes.M, sizes.M) == false {
		return false
	}
	if AreEqual(imageSizes.L, sizes.L) == false {
		return false
	}
	if AreEqual(imageSizes.XL, sizes.XL) == false {
		return false
	}
	if AreEqual(imageSizes.XXL, sizes.XXL) == false {
		return false
	}
	if AreEqual(imageSizes.XXL, sizes.XXL) == false {
		return false
	}
	if AreEqual(imageSizes.XXXL, sizes.XXXL) == false {
		return false
	}
	return true
}

func (imageSizes ImageSizes) GetAllDimensions() []Dimensions {
	dimensions := []Dimensions{imageSizes.Original}

	if imageSizes.Xs != nil {
		dimensions = append(dimensions, *imageSizes.Xs)
	}
	if imageSizes.S != nil {
		dimensions = append(dimensions, *imageSizes.S)
	}
	if imageSizes.M != nil {
		dimensions = append(dimensions, *imageSizes.M)
	}
	if imageSizes.L != nil {
		dimensions = append(dimensions, *imageSizes.L)
	}
	if imageSizes.XL != nil {
		dimensions = append(dimensions, *imageSizes.XL)
	}
	if imageSizes.XXL != nil {
		dimensions = append(dimensions, *imageSizes.XXL)
	}
	if imageSizes.XXXL != nil {
		dimensions = append(dimensions, *imageSizes.XXXL)
	}

	return dimensions
}
