package storage

import (
	"fmt"
	"time"
)

type Image struct {
	Id        string      `json:"id"`
	Name      string      `json:"name"`
	Format    ImageFormat `json:"format"`
	Original  string      `json:"original"`
	Domain    string      `json:"domain"`
	Path      string      `json:"path"`
	Sizes     ImageSizes  `json:"sizes"`
	CreatedAt *time.Time  `json:"createdAt"`
	UpdatedAt *time.Time  `json:"updatedAt"`
	AuthorId  string      `json:"authorId"`
}

func (image Image) IsEqualTo(img Image) bool {
	if image.Name != img.Name {
		return false
	}
	if image.Format != img.Format {
		return false
	}
	if image.Original != img.Original {
		return false
	}
	if image.Domain != img.Domain {
		return false
	}
	if image.Path != img.Path {
		return false
	}

	if !image.Sizes.IsEqualTo(img.Sizes) {
		return false
	}

	return true
}

func (image Image) toString() string {
	return fmt.Sprintf(
		"Image{id: %s, name: %s, format: %s, original: %s, domain: %s, path: %s, sizes: %s}",
		image.Id,
		image.Name,
		image.Format,
		image.Original,
		image.Domain,
		image.Path,
		image.Sizes.ToString(),
	)
}

type ImageList []Image

type ImagePredicate func(img Image) bool

func (images ImageList) findBy(predicate ImagePredicate) *Image {
	for _, image := range images {
		if predicate(image) {
			return &image
		}
	}
	return nil
}

func (images ImageList) IsEqualTo(imageList ImageList) bool {
	if images == nil && imageList == nil {
		return true
	}

	if images != nil && imageList != nil {
		if len(images) != len(imageList) {
			return false
		}

		for _, image := range images {
			foundImg := imageList.findBy(func(img Image) bool {
				return img.Name == image.Name
			})
			if foundImg == nil {
				return false
			}
			if !image.IsEqualTo(*foundImg) {
				return false
			}
		}
		return true
	}

	return false
}

func (images ImageList) ToString() string {
	var value string
	for _, image := range images {
		value += image.toString() + "\n"
	}

	return value
}
