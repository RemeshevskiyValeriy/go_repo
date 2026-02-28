package http

import (
	"context"
	"net/http"
	"time"

	"example.com/tech-ip-sem2/services/tasks/internal/client/authclient"
)

func AuthMiddleware(auth *authclient.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization", http.StatusUnauthorized)
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
			defer cancel()

			if err := auth.Verify(ctx, authHeader); err != nil {
				if err == authclient.ErrUnauthorized {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
				http.Error(w, "auth service unavailable", http.StatusBadGateway)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}