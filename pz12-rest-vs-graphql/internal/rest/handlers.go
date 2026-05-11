package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"example.com/pz12-rest-vs-graphql/internal/store"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.Tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	for _, t := range store.Tasks {
		if t.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
			return
		}
	}

	http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
}

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := &store.Task{
		ID:          fmt.Sprintf("t_%03d", len(store.Tasks)+1),
		Title:       req.Title,
		Description: req.Description,
		Done:        false,
	}

	store.Tasks = append(store.Tasks, task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

type UpdateTaskRequest struct {
	Done *bool `json:"done"`
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, t := range store.Tasks {
		if t.ID == id {

			if req.Done != nil {
				t.Done = *req.Done
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
			return
		}
	}

	http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
}