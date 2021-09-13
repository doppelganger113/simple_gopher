package http_transport

import (
	"context"
	"fmt"
	"mime/multipart"
	"simple_gopher/auth"
	"simple_gopher/image_resize"
	"simple_gopher/storage"
)

type ImagesHandlerMock struct {
}

func (h ImagesHandlerMock) Get(
	_ context.Context, limit, offset int, order storage.Order,
) (storage.ImageList, error) {
	if limit != 10 || offset != 10 || order != storage.OrderAscending {
		return storage.ImageList{}, fmt.Errorf(
			"expecting limit %d got %d | expecting offset %d got %d and expecting order %s got %s",
			10, limit, 10, offset, storage.OrderDescending, order,
		)
	}
	images := storage.ImageList{
		{
			Id:       "3c47d736-6c4e-4a1c-a04b-3744cc30b263",
			Name:     "my-image-1",
			Format:   "jpg",
			Original: "images/my-unique-image.png",
			Domain:   "https://random.cloudfront.net",
			Path:     "images",
			Sizes: storage.ImageSizes{
				Original: storage.Dimensions{
					Width:  688,
					Height: 516,
				},
				Xs: &storage.Dimensions{
					Width:  100,
					Height: 75,
				},
				S: &storage.Dimensions{
					Width:  300,
					Height: 225,
				},
				M: &storage.Dimensions{
					Width:  500,
					Height: 375,
				},
			},
			CreatedAt: nil,
			UpdatedAt: nil,
		},
	}

	return images, nil
}

func (h ImagesHandlerMock) GetOne(ctx context.Context, imageId string) (storage.Image, error) {
	return storage.Image{}, nil
}

func (h ImagesHandlerMock) UploadAndResize(
	ctx context.Context,
	authorization auth.AuthorizationDto,
	imageName string,
	format image_resize.ImageFormat,
	originalFile *multipart.FileHeader,
	croppedFile *multipart.FileHeader,
) (storage.Image, error) {
	return storage.Image{}, nil
}

func (h ImagesHandlerMock) DeleteOne(
	ctx context.Context,
	auth auth.AuthorizationDto,
	imageId string,
) error {
	return nil
}
