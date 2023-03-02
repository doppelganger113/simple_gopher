package core

import (
	"api/core/exception"
	"api/storage"
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (service *ImagesService) Get(
	ctx context.Context, limit, offset int, order storage.Order,
) (storage.ImageList, error) {
	images, err := service.imagesRepository.Get(ctx, limit, offset, order)
	if err != nil {
		return storage.ImageList{}, fmt.Errorf("failed fetching images: %w", err)
	}

	return images, nil
}

func (service *ImagesService) GetOne(ctx context.Context, imageId string) (storage.Image, error) {
	parsedImageId, err := uuid.Parse(imageId)
	if err != nil {
		return storage.Image{}, exception.InvalidArgument{Reason: "Invalid uuid"}
	}

	image, err := service.imagesRepository.GetOne(ctx, parsedImageId.String())
	if err != nil {
		return storage.Image{}, err
	}

	return image, nil
}
