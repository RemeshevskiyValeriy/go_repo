package httpx

import (
  "github.com/go-chi/chi/v5"
  "example.com/notes-api-swagger/internal/http/handlers"

  httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()
	
	// Swagger UI
	r.Get("/docs/*", httpSwagger.WrapHandler)
	
	r.Route("/api/v1/notes", func(r chi.Router) {
		r.Get("/", h.ListNotes)
		r.Post("/", h.CreateNote)
		r.Get("/{id}", h.GetNote)
		r.Patch("/{id}", h.PatchNote)
		r.Delete("/{id}", h.DeleteNote)
	})
	return r
}
