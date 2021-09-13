package http_transport

import (
	"context"
	"mime/multipart"
	"simple_gopher/auth"
	"simple_gopher/image_resize"
	"simple_gopher/storage"
)

type ImagesHandler interface {
	Get(ctx context.Context, limit, offset int, order storage.Order) (storage.ImageList, error)
	GetOne(ctx context.Context, imageId string) (storage.Image, error)
	UploadAndResize(
		ctx context.Context,
		authorization auth.AuthorizationDto,
		imageName string,
		format image_resize.ImageFormat,
		originalFile *multipart.FileHeader,
		croppedFile *multipart.FileHeader,
	) (storage.Image, error)
	DeleteOne(
		ctx context.Context,
		auth auth.AuthorizationDto,
		imageId string,
	) error
}
