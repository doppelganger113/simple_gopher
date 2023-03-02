package http_util

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Response struct {
	Status int
	Data   interface{}
}

func (r *Response) WithStatus(status int) *Response {
	r.Status = status
	return r
}

func (r *Response) GetStatus() int {
	if r.Status > 0 {
		return r.Status
	}

	return http.StatusOK
}

func NewResponse(data interface{}) *Response {
	return &Response{Status: 200, Data: data}
}

func WriteJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Error().Interface("data", data).Msg("error parsing err response")
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"err": "Error parsing err response"}`))
			if err != nil {
				log.Error().Err(err).Msg("error writing err response")
			}
			return
		}
	} else {
		_, _ = w.Write([]byte(""))
	}
}

func WriteBadRequestJson(w http.ResponseWriter, err error) {
	WriteJson(w, http.StatusBadRequest, NewFailureResponse(err.Error()))
}
