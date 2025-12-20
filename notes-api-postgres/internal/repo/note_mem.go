package repo

import (
  "sync"
  "example.com/notes-api-postgres/internal/core"
)

type NoteRepoMem struct {
  mu    sync.Mutex
  notes map[int64]*core.Note
  next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
  return &NoteRepoMem{notes: make(map[int64]*core.Note)}
}

func (r *NoteRepoMem) GetAll() []*core.Note {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make([]*core.Note, 0, len(r.notes))
	for _, n := range r.notes {
		res = append(res, n)
	}
	return res
}

func (r *NoteRepoMem) Create(n core.Note) (int64, error) {
  r.mu.Lock()
  defer r.mu.Unlock()

  r.next++
  n.ID = r.next
  r.notes[n.ID] = &n
  return n.ID, nil
}

func (r *NoteRepoMem) GetByID(id int64) (*core.Note, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notes[id]
	return n, ok
}

func (r *NoteRepoMem) UpdatePartial(id int64, title, content *string) (*core.Note, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notes[id]
	if !ok {
		return nil, false
	}

	if title != nil {
		n.Title = *title
	}
	if content != nil {
		n.Content = *content
	}

	return n, true
}

func (r *NoteRepoMem) Delete(id int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.notes[id]; !ok {
		return false
	}
	delete(r.notes, id)
	return true
}
