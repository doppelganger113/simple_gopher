package image_resize_api

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"simple_gopher/image_resize"
)

func (api ResizeApi) UploadFile(
	ctx context.Context,
	signedUrl string,
	format image_resize.ImageFormat,
	fileHeader *multipart.FileHeader,
) error {
	file, err := fileHeader.Open()
	contentType := string(format.ToImageContentType())

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		signedUrl,
		file,
	)
	if err != nil {
		return err
	}

	req.ContentLength = fileHeader.Size

	res, err := api.client.Do(req)
	if err != nil {
		var status int
		if res != nil {
			status = res.StatusCode
		}
		return &image_resize.BadRequest{
			RequestError: image_resize.RequestError{
				Url:        signedUrl,
				StatusCode: status,
				Message:    fmt.Sprintf("failed upload of file type %s", contentType),
				Err:        err,
			},
			Body: "",
		}
	}
	defer func() {
		if closeErr := res.Body.Close(); err != nil {
			log.Warn().Msgf("failed closing body: %s", closeErr.Error())
		}
	}()

	if isResponseOk(res.StatusCode) == false {
		body, e := ioutil.ReadAll(res.Body)
		if e != nil {
			return e
		}
		bodyMsg := fmt.Sprintf("Body: %s", body)
		fileTypeErr := fmt.Errorf("failed uploading file of type %s", contentType)

		if res.StatusCode == 403 {
			return &image_resize.Forbidden{
				RequestError: image_resize.RequestError{
					Url:        signedUrl,
					StatusCode: res.StatusCode,
					Message:    bodyMsg,
					Err:        fileTypeErr,
				},
			}
		}

		return &image_resize.BadRequest{
			RequestError: image_resize.RequestError{
				Url:        signedUrl,
				StatusCode: res.StatusCode,
				Message:    bodyMsg,
				Err:        fileTypeErr,
			},
		}
	}

	return nil
}
