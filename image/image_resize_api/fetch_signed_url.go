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

func (api *ResizeApi) FetchSignedUrl(
	ctx context.Context, authorization string, format image.Format,
) (image.SignedResponse, error) {
	request := image.SignedRequest{Format: format}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return image.SignedResponse{}, err
	}

	reqUrl := api.url("/signed")

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		reqUrl,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return image.SignedResponse{}, err
	}

	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/json")

	res, err := api.client.Do(req)
	if err != nil {
		var status int
		if res != nil {
			status = res.StatusCode
		}
		return image.SignedResponse{}, &image.BadRequest{
			RequestError: image.RequestError{
				Url:        reqUrl,
				StatusCode: status,
				Message:    "failed making signed url request",
				Err:        err,
			},
			Body: string(jsonData),
		}
	}

	defer func() {
		if closeErr := res.Body.Close(); err != nil {
			log.Warn().Msgf("error closing body: %s", closeErr.Error())
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return image.SignedResponse{}, err
	}
	if !isResponseOk(res.StatusCode) {
		return image.SignedResponse{}, &image.BadRequest{
			RequestError: image.RequestError{
				Url:        reqUrl,
				StatusCode: res.StatusCode,
				Message:    "failed reading response",
				Err:        err,
			},
			Body: string(body),
		}
	}

	var signedResponse image.SignedResponse

	err = json.Unmarshal(body, &signedResponse)
	if err != nil {
		return image.SignedResponse{}, &image.BadRequest{
			RequestError: image.RequestError{
				Url:        reqUrl,
				StatusCode: res.StatusCode,
				Message:    "failed unmarshalling response",
				Err:        err,
			},
			Body: string(body),
		}
	}

	return signedResponse, nil
}
