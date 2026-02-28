package http

import (
	"net/http"

	"example.com/tech-ip-sem2/shared/middleware"
)

func NewRouter(h *Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/auth/login", h.Login)
	mux.HandleFunc("/v1/auth/verify", h.Verify)

	// middleware chain
	handler := middleware.RequestID(mux)
	handler = middleware.Logging(handler)

	return handler
}