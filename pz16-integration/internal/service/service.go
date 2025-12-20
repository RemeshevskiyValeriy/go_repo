package service

import (
  "context"
  "example.com/pz16-integration/internal/models"
  "example.com/pz16-integration/internal/repo"
)

type Service struct{ Notes repo.NoteRepo }

func (s Service) Create(ctx context.Context, n *models.Note) error {
  // можно добавить валидацию
  return s.Notes.Create(ctx, n)
}
func (s Service) Get(ctx context.Context, id int64) (models.Note, error) {
  return s.Notes.Get(ctx, id)
}

func (s Service) Update(ctx context.Context, n *models.Note) error {
	return s.Notes.Update(ctx, n)
}

func (s Service) Delete(ctx context.Context, id int64) error {
	return s.Notes.Delete(ctx, id)
}

func (s Service) List(ctx context.Context) ([]models.Note, error) {
	return s.Notes.List(ctx)
}