package resize

import (
	"api/image"
	"context"
)

func (client *Client) Rename(
	ctx context.Context,
	authorizationHeader string,
	request image.RenameRequest,
) (image.ResizeResponse, error) {
	return image.ResizeResponse{}, nil
}
