package handlers

import (
  "encoding/json"
  "net/http"
  "strconv"
  "example.com/notes-api-swagger/internal/core"
  "example.com/notes-api-swagger/internal/repo"
  "github.com/go-chi/chi/v5"
)

type Handler struct {
  Repo *repo.NoteRepoMem
}

// ListNotes godoc
// @Summary      Список заметок
// @Description  Возвращает список заметок
// @Tags         notes
// @Success      200  {array}   core.Note
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/notes [get]
func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	notes := h.Repo.GetAll()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(notes)
}

// CreateNote godoc
// @Summary      Создать заметку
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        input  body     core.NoteCreate  true  "Данные новой заметки"
// @Success      201    {object} core.Note
// @Failure      400    {object} map[string]string
// @Failure      500    {object} map[string]string
// @Router       /api/v1/notes [post]
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
  	var input core.NoteCreate

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}

	if input.Title == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "title is required")
		return
	}

	note := core.Note{
		Title:   input.Title,
		Content: input.Content,
	}
	
	id, _ := h.Repo.Create(note)
	note.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(note)
}

// GetNote godoc
// @Summary      Получить заметку
// @Tags         notes
// @Param        id   path   int  true  "ID"
// @Success      200  {object}  core.Note
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/notes/{id} [get]
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
	_ = json.NewEncoder(w).Encode(note)
}

// PatchNote godoc
// @Summary      Обновить заметку (частично)
// @Tags         notes
// @Accept       json
// @Param        id     path   int               true  "ID"
// @Param        input  body   core.NoteUpdate   true  "Поля для обновления"
// @Success      200    {object}  core.Note
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /api/v1/notes/{id} [patch]
func (h *Handler) PatchNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "id must be integer")
		return
	}

	var input core.NoteUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}

	if input.Title == nil && input.Content == nil {
		writeError(w, http.StatusBadRequest, "validation_error", "nothing to update")
		return
	}

	note, ok := h.Repo.UpdatePartial(id, input.Title, input.Content)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "note not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(note)
}

// DeleteNote godoc
// @Summary      Удалить заметку
// @Tags         notes
// @Param        id  path  int  true  "ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/notes/{id} [delete]
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
