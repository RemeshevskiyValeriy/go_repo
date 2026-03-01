package http

import (
	"net/http"
	"strings"

	"example.com/tech-ip-grpc/services/tasks/internal/client/authgrpc"
)

func AuthMiddleware(auth *authgrpc.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			if err := auth.Verify(r.Context(), token); err != nil {
				switch err {
				case authgrpc.ErrUnauthorized:
					http.Error(w, "unauthorized", http.StatusUnauthorized)
				default:
					http.Error(w, "auth service unavailable", http.StatusServiceUnavailable)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}