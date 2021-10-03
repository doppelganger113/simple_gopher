package image

import "fmt"

type Forbidden struct {
	RequestError
	Body string
}

func (forbidden Forbidden) Error() string {
	return fmt.Sprintf(
		"{Url: %s, StatusCode: %d, Message: %s, Err: %v}",
		forbidden.Url,
		forbidden.StatusCode,
		forbidden.Message,
		forbidden.Err,
	)
}

func (forbidden Forbidden) Unwrap() error {
	return forbidden.Err
}
