package image

import "fmt"

type BadRequest struct {
	RequestError
	Body string
}

func (badRequest BadRequest) Error() string {
	return fmt.Sprintf(
		"{Url: %s, StatusCode: %d, Message: %s, Err: %v, Body: %s}",
		badRequest.Url,
		badRequest.StatusCode,
		badRequest.Message,
		badRequest.Err,
		badRequest.Body,
	)
}

func (badRequest BadRequest) Unwrap() error {
	return badRequest.Err
}
