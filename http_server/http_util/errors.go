package http_util

type FailureResponse struct {
	Err string `json:"error"`
}

func (f FailureResponse) Error() string {
	return f.Err
}

func NewFailureResponse(msg string) FailureResponse {
	return FailureResponse{Err: msg}
}

var serverErrorFailure = NewFailureResponse("Internal server error")
var forbiddenFailure = NewFailureResponse("Forbidden")
var notFoundFailure = NewFailureResponse("Not found")
