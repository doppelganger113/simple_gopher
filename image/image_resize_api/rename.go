package image_resize_api

import (
	"context"
	"simple_gopher/image"
)

func (api *ResizeApi) Rename(
	ctx context.Context,
	authorizationHeader string,
	request image.RenameRequest,
) (image.ResizeResponse, error) {
	return image.ResizeResponse{}, nil
}
