package image_resize

import (
	"context"
	"mime/multipart"
)

type Resizer interface {
	FetchSignedUrl(
		ctx context.Context,
		authorization string,
		format ImageFormat,
	) (*SignedResponse, error)
	UploadFile(
		ctx context.Context,
		signedUrl string,
		format ImageFormat,
		fileHeader *multipart.FileHeader,
	) error
	Resize(
		ctx context.Context,
		authorizationHeader string,
		imageResizeRequest ImageResizeRequest,
	) (ImageResizeResponse, error)
	Rename(
		ctx context.Context,
		authorizationHeader string,
		request ImageRenameRequest,
	) (ImageResizeResponse, error)
	Invalidate(
		ctx context.Context,
		authorizationHeader string,
		request ImageDeleteRequest,
	) error
	Delete(
		ctx context.Context,
		authorizationHeader string,
		request ImageDeleteRequest,
	) error
}
