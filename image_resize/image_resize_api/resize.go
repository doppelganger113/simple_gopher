package image_resize_api

import (
	"bytes"
	"context"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"simple_gopher/image_resize"
)

func (api *ResizeApi) Resize(
	ctx context.Context,
	authorizationHeader string,
	imageResizeRequest image_resize.ImageResizeRequest,
) (image_resize.ImageResizeResponse, error) {

	jsonData, err := json.Marshal(imageResizeRequest)
	if err != nil {
		return image_resize.ImageResizeResponse{}, err
	}

	requestUrl := api.url("/resize")

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		requestUrl,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return image_resize.ImageResizeResponse{}, err
	}

	req.Header.Add("Authorization", authorizationHeader)
	req.Header.Add("Content-Type", "application/json")

	log.Info().Caller().Msg("issuing resize request")
	res, err := api.client.Do(req)
	if err != nil {
		var statusCode int
		if res != nil {
			statusCode = res.StatusCode
		}
		return image_resize.ImageResizeResponse{}, &image_resize.BadRequest{
			RequestError: image_resize.RequestError{
				Url:        requestUrl,
				StatusCode: statusCode,
				Message:    "Failed making resize request",
				Err:        err,
			},
		}
	}
	defer func() {
		if closeErr := res.Body.Close(); err != nil {
			log.Warn().Caller().Msgf("error closing %s", closeErr.Error())
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return image_resize.ImageResizeResponse{}, err
	}
	if isResponseOk(res.StatusCode) == false {
		if res.StatusCode == 403 {
			return image_resize.ImageResizeResponse{}, &image_resize.Forbidden{
				RequestError: image_resize.RequestError{
					Url:        requestUrl,
					StatusCode: res.StatusCode,
					Message:    "Forbidden request",
				},
				Body: string(body),
			}
		}
		return image_resize.ImageResizeResponse{}, &image_resize.BadRequest{
			RequestError: image_resize.RequestError{
				Url:        requestUrl,
				StatusCode: res.StatusCode,
				Message:    "Failed resize request",
			},
			Body: string(body),
		}
	}

	var response image_resize.ImageResizeResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return image_resize.ImageResizeResponse{}, err
	}

	return response, nil
}
