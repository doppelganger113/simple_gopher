package image

import (
	"context"
	"mime/multipart"
)

type Mock struct {
}

func (resize Mock) FetchSignedUrl(
	ctx context.Context,
	authorization string,
	format Format,
) (SignedResponse, error) {
	return SignedResponse{}, nil
}

func (resize Mock) UploadFile(
	ctx context.Context,
	signedUrl string,
	format Format,
	fileHeader *multipart.FileHeader,
) error {
	return nil
}

func (resize Mock) Resize(
	ctx context.Context, authorizationHeader string, imageResizeRequest ResizeRequest,
) (ResizeResponse, error) {
	return ResizeResponse{}, nil
}

func (resize Mock) Invalidate(
	ctx context.Context,
	authorizationHeader string,
	request DeleteRequest,
) error {
	return nil
}

func (resize Mock) Rename(
	ctx context.Context,
	authorizationHeader string,
	request RenameRequest,
) (ResizeResponse, error) {
	return ResizeResponse{}, nil
}
func (resize Mock) Delete(
	ctx context.Context,
	authorizationHeader string,
	request DeleteRequest,
) error {
	return nil
}
