package main

import (
  "log"
  "net/http"
  "example.com/notes-api-swagger/internal/http"
  "example.com/notes-api-swagger/internal/http/handlers"
  "example.com/notes-api-swagger/internal/repo"

	_ "example.com/notes-api-swagger/docs"
)

// Package main Notes API server.
//
// @title           Notes API
// @version         1.0
// @description     Учебный REST API для заметок (CRUD).
// @contact.name    Backend Course
// @contact.email   example@university.ru
// @BasePath        /api/v1
func main() {
  repo := repo.NewNoteRepoMem()
  h := &handlers.Handler{Repo: repo}
  r := httpx.NewRouter(h)

  log.Println("Server started at :8080")
  log.Fatal(http.ListenAndServe(":8080", r))
}
