package resize

import (
	"api/image"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func (client *Client) UploadFile(
	ctx context.Context,
	signedUrl string,
	format image.Format,
	fileHeader *multipart.FileHeader,
) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed opening file header: %w", err)
	}
	contentType := string(format.ToContentType())

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

	res, err := client.client.Do(req)
	if err != nil {
		var status int
		if res != nil {
			status = res.StatusCode
		}
		return &image.BadRequest{
			RequestError: image.RequestError{
				Url:        signedUrl,
				StatusCode: status,
				Message:    fmt.Sprintf("failed upload of file type %s", contentType),
				Err:        err,
			},
			Body: "",
		}
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			client.logger.Warn().Msgf("failed closing body: %s", closeErr.Error())
		}
	}()

	if !isResponseOk(res.StatusCode) {
		body, e := ioutil.ReadAll(res.Body)
		if e != nil {
			return e
		}
		bodyMsg := fmt.Sprintf("Body: %s", body)
		fileTypeErr := fmt.Errorf("failed uploading file of type %s", contentType)

		if res.StatusCode == 403 {
			return &image.Forbidden{
				RequestError: image.RequestError{
					Url:        signedUrl,
					StatusCode: res.StatusCode,
					Message:    bodyMsg,
					Err:        fileTypeErr,
				},
			}
		}

		return &image.BadRequest{
			RequestError: image.RequestError{
				Url:        signedUrl,
				StatusCode: res.StatusCode,
				Message:    bodyMsg,
				Err:        fileTypeErr,
			},
		}
	}

	return nil
}
