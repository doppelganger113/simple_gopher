package concurrency

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type ContextKey string

var (
	CorrelationIdKey ContextKey = "X-Correlation-Id"
)

// CorrelationId middleware adds correlation id to the context from the header if it exists or creates one
// Access is it through requests
//  val := req.Context().Value(concurrency.CorrelationIdKey)
func CorrelationId(next http.HandlerFunc, errHandler func(err error)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		correlationId := req.Header.Get(string(CorrelationIdKey))
		if correlationId == "" {
			id, err := uuid.NewUUID()
			if err != nil {
				errHandler(err)
				next(w, req)
				return
			}

			correlationId = id.String()
		}
		updatedCtx := context.WithValue(req.Context(), CorrelationIdKey, correlationId)
		updatedReq := req.WithContext(updatedCtx)

		next(w, updatedReq)
	}
}
