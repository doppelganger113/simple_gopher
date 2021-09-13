package image_resize

import (
	"context"
	"mime/multipart"
)

type Mock struct {
}

func (resize Mock) FetchSignedUrl(
	ctx context.Context,
	authorization string,
	format ImageFormat,
) (SignedResponse, error) {
	return SignedResponse{}, nil
}

func (resize Mock) UploadFile(
	ctx context.Context,
	signedUrl string,
	format ImageFormat,
	fileHeader *multipart.FileHeader,
) error {
	return nil
}

func (resize Mock) Resize(
	ctx context.Context, authorizationHeader string, imageResizeRequest ImageResizeRequest,
) (ImageResizeResponse, error) {
	return ImageResizeResponse{}, nil
}

func (resize Mock) Invalidate(
	ctx context.Context,
	authorizationHeader string,
	request ImageDeleteRequest,
) error {
	return nil
}

func (resize Mock) Rename(
	ctx context.Context,
	authorizationHeader string,
	request ImageRenameRequest,
) (ImageResizeResponse, error) {
	return ImageResizeResponse{}, nil
}
func (resize Mock) Delete(
	ctx context.Context,
	authorizationHeader string,
	request ImageDeleteRequest,
) error {
	return nil
}
