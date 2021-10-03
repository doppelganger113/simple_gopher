package image

import "fmt"

type RequestError struct {
	Url        string
	StatusCode int
	Message    string
	Err        error
}

func (err RequestError) Error() string {
	return fmt.Sprintf(
		"{Url: %s, StatusCode: %d, Message: %s, Err: %v}",
		err.Url,
		err.StatusCode,
		err.Message,
		err.Err,
	)
}

func (err RequestError) Unwrap() error {
	return err.Err
}
