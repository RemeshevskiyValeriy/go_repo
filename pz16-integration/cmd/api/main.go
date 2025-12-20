// cmd/api/main.go
package main

import (
	"log"
	"os"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"example.com/pz16-integration/internal/db"
	"example.com/pz16-integration/internal/httpapi"
	"example.com/pz16-integration/internal/repo"
	"example.com/pz16-integration/internal/service"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}

	dbx, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.MustApplyMigrations(dbx)

	svc := service.Service{
		Notes: repo.NoteRepo{DB: dbx},
	}

	r := gin.Default()
	httpapi.Router{Svc: &svc}.Register(r)

	log.Println("API listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
