package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"

	"example.com/notes-api-postgres/internal/core"
)

type NoteRepoPG struct {
	db *sql.DB
}

func NewNoteRepoPG(db *sql.DB) *NoteRepoPG {
	return &NoteRepoPG{db: db}
}

func (r *NoteRepoPG) ListOffset(
	ctx context.Context,
	offset, limit int,
) ([]*core.Note, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, content, created_at, updated_at
		FROM notes
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2
	`, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*core.Note
	for rows.Next() {
		var n core.Note
		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Content,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		notes = append(notes, &n)
	}

	return notes, rows.Err()
}

func (r *NoteRepoPG) GetByID(
	ctx context.Context,
	id int64,
) (*core.Note, error) {

	var n core.Note
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, content, created_at, updated_at
		FROM notes
		WHERE id = $1
	`, id).Scan(
		&n.ID,
		&n.Title,
		&n.Content,
		&n.CreatedAt,
		&n.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *NoteRepoPG) Create(
	ctx context.Context,
	n core.Note,
) (int64, error) {

	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO notes (title, content)
		VALUES ($1, $2)
		RETURNING id
	`, n.Title, n.Content).Scan(&id)

	return id, err
}

func (r *NoteRepoPG) ListKeyset(
	ctx context.Context,
	lastCreatedAt time.Time,
	lastID int64,
	limit int,
) ([]*core.Note, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, content, created_at, updated_at
		FROM notes
		WHERE (created_at, id) < ($1, $2)
		ORDER BY created_at DESC, id DESC
		LIMIT $3
	`, lastCreatedAt, lastID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*core.Note
	for rows.Next() {
		var n core.Note
		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Content,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		notes = append(notes, &n)
	}

	return notes, rows.Err()
}

func (r *NoteRepoPG) GetByIDs(
	ctx context.Context,
	ids []int64,
) ([]*core.Note, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, content, created_at, updated_at
		FROM notes
		WHERE id = ANY($1)
	`, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*core.Note
	for rows.Next() {
		var n core.Note
		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Content,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		notes = append(notes, &n)
	}

	return notes, rows.Err()
}


