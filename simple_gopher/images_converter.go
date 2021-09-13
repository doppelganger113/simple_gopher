package simple_gopher

import (
	"simple_gopher/image_resize"
	"simple_gopher/storage"
)

func toImageResizeDimensions(dimensions storage.Dimensions) image_resize.Dimensions {
	return image_resize.Dimensions{
		Width:  dimensions.Width,
		Height: dimensions.Height,
	}
}

func convertStorageSizesToDimensions(imageSizes storage.ImageSizes) []image_resize.Dimensions {
	dimensions := []image_resize.Dimensions{toImageResizeDimensions(imageSizes.Original)}

	if imageSizes.Xs != nil {
		dimensions = append(dimensions, toImageResizeDimensions(*imageSizes.Xs))
	}
	if imageSizes.S != nil {
		dimensions = append(dimensions, toImageResizeDimensions(*imageSizes.S))
	}
	if imageSizes.M != nil {
		dimensions = append(dimensions, toImageResizeDimensions(*imageSizes.M))
	}
	if imageSizes.L != nil {
		dimensions = append(dimensions, toImageResizeDimensions(*imageSizes.L))
	}
	if imageSizes.XL != nil {
		dimensions = append(dimensions, toImageResizeDimensions(*imageSizes.XL))
	}
	if imageSizes.XXL != nil {
		dimensions = append(dimensions, toImageResizeDimensions(*imageSizes.XXL))
	}
	if imageSizes.XXXL != nil {
		dimensions = append(dimensions, toImageResizeDimensions(*imageSizes.XXXL))
	}

	return dimensions
}

func fromImageResizeDimensions(img *image_resize.Dimensions) *storage.Dimensions {
	return &storage.Dimensions{
		Width:  img.Width,
		Height: img.Height,
	}
}

func fromStorageDimensions(img *storage.Dimensions) *image_resize.Dimensions {
	return &image_resize.Dimensions{
		Width:  img.Width,
		Height: img.Height,
	}
}

func convertImageSizesToStorageSizes(sizes image_resize.ImageSizes) storage.ImageSizes {
	return storage.ImageSizes{
		Original: storage.Dimensions{
			Width:  sizes.Original.Width,
			Height: sizes.Original.Height,
		},
		Xs:   fromImageResizeDimensions(sizes.Xs),
		S:    fromImageResizeDimensions(sizes.S),
		M:    fromImageResizeDimensions(sizes.M),
		L:    fromImageResizeDimensions(sizes.L),
		XL:   fromImageResizeDimensions(sizes.XL),
		XXL:  fromImageResizeDimensions(sizes.XXL),
		XXXL: fromImageResizeDimensions(sizes.XXXL),
	}
}

func fromStorageImageSizesToImageSizes(sizes storage.ImageSizes) image_resize.ImageSizes {
	return image_resize.ImageSizes{
		Original: image_resize.Dimensions{
			Width:  sizes.Original.Width,
			Height: sizes.Original.Height,
		},
		Xs:   fromStorageDimensions(sizes.Xs),
		S:    fromStorageDimensions(sizes.S),
		M:    fromStorageDimensions(sizes.M),
		L:    fromStorageDimensions(sizes.L),
		XL:   fromStorageDimensions(sizes.XL),
		XXL:  fromStorageDimensions(sizes.XXL),
		XXXL: fromStorageDimensions(sizes.XXXL),
	}
}
