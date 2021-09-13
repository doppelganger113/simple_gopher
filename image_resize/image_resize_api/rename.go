package image_resize_api

import (
	"context"
	"simple_gopher/image_resize"
)

func (api ResizeApi) Rename(
	ctx context.Context,
	authorizationHeader string,
	request image_resize.ImageRenameRequest,
) (image_resize.ImageResizeResponse, error) {
	return image_resize.ImageResizeResponse{}, nil
}
