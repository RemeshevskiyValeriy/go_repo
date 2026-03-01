package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"example.com/tech-ip-grpc/services/tasks/internal/service"
)

type Handlers struct {
	tasks *service.TaskService
}

func NewHandlers(tasks *service.TaskService) *Handlers {
	return &Handlers{tasks: tasks}
}

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var t service.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if t.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	created := h.tasks.Create(&t)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.tasks.List())
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	t, err := h.tasks.Get(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(t)
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	var upd service.Task
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	t, err := h.tasks.Update(id, &upd)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(t)
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	if err := h.tasks.Delete(id); err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}