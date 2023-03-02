package http_util

import (
	"context"
	"github.com/rs/zerolog"
	"net/http"
)

type RequestHandler struct {
	logger *zerolog.Logger
}

func NewRequestHandler(logger *zerolog.Logger) RequestHandler {
	return RequestHandler{logger: logger}
}

type BasicRequestHandler func(ctx context.Context, req *http.Request) (*Response, error)

func (h RequestHandler) Handle(fn BasicRequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		response, err := fn(ctx, req)
		if err != nil {
			HandleError(h.logger, w, err)
		} else {
			WriteJson(w, response.GetStatus(), response.Data)
		}
	}
}
