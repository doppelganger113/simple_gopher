package image

import (
	"context"
	"mime/multipart"
)

type Resizer interface {
	FetchSignedUrl(
		ctx context.Context,
		authorization string,
		format Format,
	) (SignedResponse, error)
	UploadFile(
		ctx context.Context,
		signedUrl string,
		format Format,
		fileHeader *multipart.FileHeader,
	) error
	Resize(
		ctx context.Context,
		authorizationHeader string,
		imageResizeRequest ResizeRequest,
	) (ResizeResponse, error)
	Rename(
		ctx context.Context,
		authorizationHeader string,
		request RenameRequest,
	) (ResizeResponse, error)
	Invalidate(
		ctx context.Context,
		authorizationHeader string,
		request DeleteRequest,
	) error
	Delete(
		ctx context.Context,
		authorizationHeader string,
		request DeleteRequest,
	) error
}
