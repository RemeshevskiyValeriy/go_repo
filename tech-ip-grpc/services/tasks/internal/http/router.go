package http

import (
	"net/http"

	"example.com/tech-ip-grpc/shared/middleware"
	"example.com/tech-ip-grpc/services/tasks/internal/client/authgrpc"
)

func NewRouter(h *Handlers, auth *authgrpc.Client) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodGet:
			h.List(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.Get(w, r)
		case http.MethodPatch:
			h.Update(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	handler := AuthMiddleware(auth)(mux)
	handler = middleware.RequestID(handler)
	handler = middleware.Logging(handler)

	return handler
}