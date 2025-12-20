package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

  "github.com/joho/godotenv"

	"example.com/notes-api-postgres/internal/http"
	"example.com/notes-api-postgres/internal/http/handlers"
	"example.com/notes-api-postgres/internal/repo"

	_ "example.com/notes-api-postgres/docs"
  _ "github.com/jackc/pgx/v5/stdlib"
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
  if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
  
  db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
  if err != nil { log.Fatal(err) }

  db.SetMaxOpenConns(5)
  db.SetMaxIdleConns(2)
  db.SetConnMaxLifetime(30 * time.Minute)
  db.SetConnMaxIdleTime(5 * time.Minute)

  
  if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

  noteRepo := repo.NewNoteRepoPG(db)

  h := &handlers.Handler{
    Repo: noteRepo,
  }

  r := httpx.NewRouter(h)

  log.Println("Server started at :8080")
  log.Fatal(http.ListenAndServe(":8080", r))
}
