package image_resize_api

import (
	"bytes"
	"context"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"net/http"
	"simple_gopher/image"
)

func (api *ResizeApi) Delete(
	ctx context.Context,
	authorizationHeader string,
	request image.DeleteRequest,
) error {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}
	reqUrl := api.url("/delete")

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		reqUrl,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", authorizationHeader)
	req.Header.Add("Content-Type", "application/json")

	res, err := api.client.Do(req)
	if err != nil {
		var status int
		if res != nil {
			status = res.StatusCode
		}
		return &image.BadRequest{
			RequestError: image.RequestError{
				Url:        reqUrl,
				StatusCode: status,
				Message:    "failed deleting",
				Err:        err,
			},
			Body: string(jsonData),
		}
	}
	if err = res.Body.Close(); err != nil {
		log.Warn().Msgf("error closing body: %s", err.Error())
	}

	return nil
}
