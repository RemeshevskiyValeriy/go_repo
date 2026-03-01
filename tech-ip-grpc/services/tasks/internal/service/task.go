package service

import (
	"errors"
	"fmt"
	"sync"
)

var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Done        bool   `json:"done"`
}

type TaskService struct {
	mu    sync.Mutex
	tasks map[string]*Task
	next  int
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks: make(map[string]*Task),
		next:  1,
	}
}

func (s *TaskService) Create(task *Task) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.next
	s.next++

	task.ID = "t_" + formatID(id)
	s.tasks[task.ID] = task
	return task
}

func (s *TaskService) List() []*Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		res = append(res, t)
	}
	return res
}

func (s *TaskService) Get(id string) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}
	return t, nil
}

func (s *TaskService) Update(id string, upd *Task) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}

	if upd.Title != "" {
		t.Title = upd.Title
	}
	if upd.Description != "" {
		t.Description = upd.Description
	}
	t.Done = upd.Done

	return t, nil
}

func (s *TaskService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrTaskNotFound
	}
	delete(s.tasks, id)
	return nil
}

func formatID(n int) string {
	if n < 10 {
		return "00" + itoa(n)
	}
	if n < 100 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}