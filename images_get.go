package simple_gopher

import (
	"context"
	"github.com/google/uuid"
	"simple_gopher/storage"
)

func (service *ImagesService) Get(
	ctx context.Context, limit, offset int, order storage.Order,
) (storage.ImageList, error) {
	return service.imagesRepository.Get(ctx, limit, offset, order)
}

func (service *ImagesService) GetOne(ctx context.Context, imageId string) (*storage.Image, error) {
	parsedImageId, err := uuid.Parse(imageId)
	if err != nil {
		return nil, InvalidArgument{Reason: "Invalid uuid"}
	}

	image, err := service.imagesRepository.GetOne(ctx, parsedImageId.String())
	if err != nil {
		return nil, err
	}

	return image, nil
}
