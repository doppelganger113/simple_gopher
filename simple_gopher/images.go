package simple_gopher

import (
	"simple_gopher/auth"
	"simple_gopher/image_resize"
	"simple_gopher/storage"
)

type ImagesService struct {
	resizeApi        image_resize.Resizer
	imagesRepository storage.ImagesRepository
	authenticator    auth.Authenticator
}

func NewImagesService(
	resizeApi image_resize.Resizer,
	imagesRepository storage.ImagesRepository,
	authenticator auth.Authenticator,
) ImagesService {
	return ImagesService{
		resizeApi:        resizeApi,
		imagesRepository: imagesRepository,
		authenticator:    authenticator,
	}
}
