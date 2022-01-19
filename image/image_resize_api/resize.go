package image_resize_api

import (
	"bytes"
	"context"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"simple_gopher/image"
)

func (api *ResizeApi) Resize(
	ctx context.Context,
	authorizationHeader string,
	imageResizeRequest image.ResizeRequest,
) (image.ResizeResponse, error) {

	jsonData, err := json.Marshal(imageResizeRequest)
	if err != nil {
		return image.ResizeResponse{}, err
	}

	requestUrl := api.url("/resize")

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		requestUrl,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return image.ResizeResponse{}, err
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
		return image.ResizeResponse{}, &image.BadRequest{
			RequestError: image.RequestError{
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
		return image.ResizeResponse{}, err
	}
	if !isResponseOk(res.StatusCode) {
		if res.StatusCode == 403 {
			return image.ResizeResponse{}, &image.Forbidden{
				RequestError: image.RequestError{
					Url:        requestUrl,
					StatusCode: res.StatusCode,
					Message:    "Forbidden request",
				},
				Body: string(body),
			}
		}
		return image.ResizeResponse{}, &image.BadRequest{
			RequestError: image.RequestError{
				Url:        requestUrl,
				StatusCode: res.StatusCode,
				Message:    "Failed resize request",
			},
			Body: string(body),
		}
	}

	var response image.ResizeResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return image.ResizeResponse{}, err
	}

	return response, nil
}
