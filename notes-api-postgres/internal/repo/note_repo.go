package repo

import (
	"context"
	"time"

	"example.com/notes-api-postgres/internal/core"
)

type NoteRepo interface {
	ListOffset(ctx context.Context, offset, limit int) ([]*core.Note, error)
	GetByID(ctx context.Context, id int64) (*core.Note, error)
	Create(ctx context.Context, n core.Note) (int64, error)
	ListKeyset(
		ctx context.Context,
		lastCreatedAt time.Time,
		lastID int64,
		limit int,
	) ([]*core.Note, error)
	GetByIDs(ctx context.Context, ids []int64) ([]*core.Note, error)
}
