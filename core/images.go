package core

import (
	"api/auth"
	"api/image"
	"api/storage"
	"github.com/rs/zerolog"
)

type ImagesService struct {
	resizeApi        image.Resizer
	imagesRepository storage.ImagesRepository
	authenticator    auth.Authenticator
	logger           *zerolog.Logger
}

func NewImagesService(
	resizeApi image.Resizer,
	imagesRepository storage.ImagesRepository,
	authenticator auth.Authenticator,
	logger *zerolog.Logger,
) *ImagesService {
	return &ImagesService{
		resizeApi:        resizeApi,
		imagesRepository: imagesRepository,
		authenticator:    authenticator,
		logger:           logger,
	}
}
