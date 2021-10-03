package simple_gopher

import (
	"simple_gopher/auth"
	"simple_gopher/image"
	"simple_gopher/storage"
)

type ImagesService struct {
	resizeApi        image.Resizer
	imagesRepository storage.ImagesRepository
	authenticator    auth.Authenticator
}

func NewImagesService(
	resizeApi image.Resizer,
	imagesRepository storage.ImagesRepository,
	authenticator auth.Authenticator,
) *ImagesService {
	return &ImagesService{
		resizeApi:        resizeApi,
		imagesRepository: imagesRepository,
		authenticator:    authenticator,
	}
}
