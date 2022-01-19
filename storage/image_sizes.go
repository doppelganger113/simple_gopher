package storage

import "fmt"

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
	if !AreEqual(&imageSizes.Original, &sizes.Original) {
		return false
	}
	if !AreEqual(imageSizes.Xs, sizes.Xs) {
		return false
	}
	if !AreEqual(imageSizes.S, sizes.S) {
		return false
	}
	if !AreEqual(imageSizes.M, sizes.M) {
		return false
	}
	if !AreEqual(imageSizes.L, sizes.L) {
		return false
	}
	if !AreEqual(imageSizes.XL, sizes.XL) {
		return false
	}
	if !AreEqual(imageSizes.XXL, sizes.XXL) {
		return false
	}
	if !AreEqual(imageSizes.XXL, sizes.XXL) {
		return false
	}
	if !AreEqual(imageSizes.XXXL, sizes.XXXL) {
		return false
	}
	return true
}
