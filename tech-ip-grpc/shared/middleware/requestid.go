package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestIDKeyType struct{}

var requestIDKey = requestIDKeyType{}

const HeaderRequestID = "X-Request-ID"

// GetRequestID — достаёт request-id из контекста
func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		if rid, ok := v.(string); ok {
			return rid
		}
	}
	return ""
}

// RequestID middleware
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(HeaderRequestID)
		if rid == "" {
			rid = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, rid)
		r = r.WithContext(ctx)

		w.Header().Set(HeaderRequestID, rid)

		next.ServeHTTP(w, r)
	})
}