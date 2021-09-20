package cloud_patterns

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const CorrelationIdKey = "X-Correlation-Id"

// CorrelationId middleware adds correlation id to the context from the header if it exists or creates one
// Access is it through requests
//  val := req.Context().Value(cloud_patterns.CorrelationIdKey)
func CorrelationId(next http.HandlerFunc, errHandler func(err error)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		correlationId := req.Header.Get(CorrelationIdKey)
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
