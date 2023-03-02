package resize

import (
	"api/image"
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (client *Client) Resize(
	ctx context.Context,
	authorizationHeader string,
	imageResizeRequest image.ResizeRequest,
) (image.ResizeResponse, error) {

	jsonData, err := json.Marshal(imageResizeRequest)
	if err != nil {
		return image.ResizeResponse{}, err
	}

	requestUrl := client.url("/resize")

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

	client.logger.Info().Msg("issuing resize request")

	res, err := client.client.Do(req)
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
		if closeErr := res.Body.Close(); closeErr != nil {
			client.logger.Warn().Msgf("failed closing body: %s", closeErr.Error())
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
