package handlers

import (
  "encoding/json"
  "net/http"
  "strconv"
  "example.com/notes-api/internal/core"
  "example.com/notes-api/internal/repo"
  "github.com/go-chi/chi/v5"
)

type Handler struct {
  Repo *repo.NoteRepoMem
}

type NotePatch struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	notes := h.Repo.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
  var n core.Note

	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}

	if n.Title == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "title is required")
		return
	}

	id, _ := h.Repo.Create(n)
	n.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}

func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "id must be integer")
		return
	}

	note, ok := h.Repo.GetByID(id)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "note not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) PatchNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "id must be integer")
		return
	}

	var patch NotePatch
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}

	if patch.Title == nil && patch.Content == nil {
		writeError(w, http.StatusBadRequest, "validation_error", "nothing to update")
		return
	}

	note, ok := h.Repo.UpdatePartial(id, patch.Title, patch.Content)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "note not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "id must be integer")
		return
	}

	if !h.Repo.Delete(id) {
		writeError(w, http.StatusNotFound, "not_found", "note not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
