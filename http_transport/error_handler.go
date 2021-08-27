package http_transport

import (
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"simple_gopher"
)

type FailureResponse struct {
	Err string `json:"err"`
}

func newFailureResponse(msg string) *FailureResponse {
	return &FailureResponse{Err: msg}
}

var serverErrorFailure = newFailureResponse("Server error")
var forbiddenFailure = newFailureResponse("Forbidden")
var notFoundFailure = newFailureResponse("Not found")

func handleError(w http.ResponseWriter, err error) {

	var forbiddenFail simple_gopher.Forbidden
	if errors.As(err, &forbiddenFail) {
		respondJson(w, http.StatusForbidden, forbiddenFailure)
		return
	}

	var notFoundFail simple_gopher.NotFound
	if errors.As(err, &notFoundFail) {
		respondJson(w, http.StatusNotFound, notFoundFailure)
		return
	}

	var invalidArgumentFail simple_gopher.InvalidArgument
	if errors.As(err, &invalidArgumentFail) {
		respondJson(w, http.StatusBadRequest, &FailureResponse{
			Err: invalidArgumentFail.Error(),
		})
		return
	}

	log.Err(err).Msg("Internal Server Error")

	respondJson(w, http.StatusInternalServerError, serverErrorFailure)
}
