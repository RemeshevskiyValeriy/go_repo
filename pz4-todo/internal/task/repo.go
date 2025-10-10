package task

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

var ErrNotFound = errors.New("task not found")

type Repo struct {
	mu        sync.RWMutex
	seq       int64
	items     map[int64]*Task
	file_path string
}

func NewRepo(file_path string) *Repo {
	r := &Repo{
		items:    make(map[int64]*Task),
		file_path: file_path,
	}
	r.load()
	return r
}

func (r *Repo) load() {
	data, err := os.ReadFile(r.file_path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}

	var tasks []*Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		panic(err)
	}

	r.items = make(map[int64]*Task)
	for _, t := range tasks {
		r.items[t.ID] = t
		if t.ID > r.seq {
			r.seq = t.ID
		}
	}
}

func (r *Repo) save() {
	out := make([]*Task, 0, len(r.items))
	for _, t := range r.items {
		out = append(out, t)
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(r.file_path, data, 0644)
}

func (r *Repo) List() []*Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Task, 0, len(r.items))
	for _, t := range r.items {
		out = append(out, t)
	}
	return out
}

func (r *Repo) Get(id int64) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (r *Repo) Create(title string) *Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	now := time.Now()
	t := &Task{ID: r.seq, Title: title, CreatedAt: now, UpdatedAt: now, Done: false}
	r.items[t.ID] = t
	r.save()
	return t
}

func (r *Repo) Update(id int64, title string, done bool) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	t.Title = title
	t.Done = done
	t.UpdatedAt = time.Now()
	r.save()
	return t, nil
}

func (r *Repo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	r.save()
	return nil
}
