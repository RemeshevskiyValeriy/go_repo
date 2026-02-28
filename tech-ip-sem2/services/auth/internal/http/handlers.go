package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"example.com/tech-ip-sem2/services/auth/internal/service"
)

type Handlers struct {
	auth *service.AuthService
}

func NewHandlers(auth *service.AuthService) *Handlers {
	return &Handlers{auth: auth}
}

/*
POST /v1/auth/login
*/
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	token, err := h.auth.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	resp := map[string]string{
		"access_token": token,
		"token_type":   "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/*
GET /v1/auth/verify
*/
func (h *Handlers) Verify(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "invalid authorization format", http.StatusUnauthorized)
		return
	}

	token := parts[1]

	subject, err := h.auth.VerifyToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"valid": false,
			"error": "unauthorized",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"valid":   true,
		"subject": subject,
	})
}