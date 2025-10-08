package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/pz3-http/internal/storage"
)

func setupHandlers() *Handlers {
	store := storage.NewMemoryStore()
	return NewHandlers(store)
}

func TestHealthHandler(t *testing.T) {
	httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Прямой вызов, как в main.go
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %v", body["status"])
	}
}

func TestCreateTask(t *testing.T) {
	h := setupHandlers()
	body := `{"title":"Test Task"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}
	var task map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&task)
	if task["title"] != "Test Task" {
		t.Errorf("expected title 'Test Task', got %v", task["title"])
	}
}

func TestCreateTask_TitleLength(t *testing.T) {
	h := setupHandlers()

	// Title слишком короткий
	body := `{"title":"ab"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateTask(w, req)
	resp := w.Result()
	if resp.StatusCode != 422 {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}

	// Title слишком длинный
	longTitle := make([]byte, 141)
	for i := range longTitle {
		longTitle[i] = 'a'
	}
	body = `{"title":"` + string(longTitle) + `"}`
	req = httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	h.CreateTask(w, req)
	resp = w.Result()
	if resp.StatusCode != 422 {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}
}

func TestListTasks(t *testing.T) {
	h := setupHandlers()
	// Добавим задачу
	h.Store.Create("Task 1")
	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	h.ListTasks(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var tasks []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

func TestGetTask(t *testing.T) {
	h := setupHandlers()
	task := h.Store.Create("Task 2")
	req := httptest.NewRequest("GET", "/tasks/"+fmt.Sprintf("%d", task.ID), nil)
	w := httptest.NewRecorder()

	h.GetTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var got map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&got)
	if got["title"] != "Task 2" {
		t.Errorf("expected title 'Task 2', got %v", got["title"])
	}
}

func TestMarkTaskDone(t *testing.T) {
	h := setupHandlers()
	h.Store.Create("Task 3")
	req := httptest.NewRequest("PATCH", "/tasks/1", nil)
	w := httptest.NewRecorder()

	h.MarkTaskDone(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var got map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&got)
	if got["done"] != true {
		t.Errorf("expected done true, got %v", got["done"])
	}
}

func TestDeleteTask(t *testing.T) {
	h := setupHandlers()
	h.Store.Create("Task 4")
	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	w := httptest.NewRecorder()

	h.DeleteTask(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}
}
