package image_resize

import (
	"context"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
)

type Mock struct {
	mock.Mock
}

func (resize *Mock) FetchSignedUrl(
	ctx context.Context,
	authorization string,
	format ImageFormat,
) (*SignedResponse, error) {
	args := resize.Called(ctx, authorization, format)

	return args.Get(0).(*SignedResponse), args.Error(1)
}

func (resize *Mock) UploadFile(
	ctx context.Context,
	signedUrl string,
	format ImageFormat,
	fileHeader *multipart.FileHeader,
) error {
	args := resize.Called(ctx, signedUrl, format, fileHeader)

	return args.Error(0)
}

func (resize *Mock) Resize(
	ctx context.Context, authorizationHeader string, imageResizeRequest ImageResizeRequest,
) (ImageResizeResponse, error) {
	args := resize.Called(ctx, authorizationHeader, imageResizeRequest)

	return args.Get(0).(ImageResizeResponse), args.Error(1)
}

func (resize *Mock) Invalidate(
	ctx context.Context,
	authorizationHeader string,
	request ImageDeleteRequest,
) error {
	args := resize.Called(ctx, authorizationHeader, request)

	return args.Error(0)
}

func (resize *Mock) Rename(
	ctx context.Context,
	authorizationHeader string,
	request ImageRenameRequest,
) (ImageResizeResponse, error) {
	args := resize.Called(ctx, authorizationHeader, request)

	return args.Get(0).(ImageResizeResponse), args.Error(1)
}
func (resize *Mock) Delete(
	ctx context.Context,
	authorizationHeader string,
	request ImageDeleteRequest,
) error {
	args := resize.Called(ctx, authorizationHeader, request)

	return args.Error(0)
}
