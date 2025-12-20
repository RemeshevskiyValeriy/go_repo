package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"example.com/pz16-integration/internal/db"
	"example.com/pz16-integration/internal/httpapi"
	"example.com/pz16-integration/internal/repo"
	"example.com/pz16-integration/internal/service"
)

func withPostgres(t *testing.T) (dsn string, term func()) {
	t.Helper()

	ctx := context.Background()
	pg, err := postgres.RunContainer(ctx,
		postgres.WithDatabase("notes_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		t.Fatal(err)
	}

	host, err := pg.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	port, err := pg.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	dsn = fmt.Sprintf(
		"postgres://test:test@%s:%s/notes_test?sslmode=disable",
		host,
		port.Port(),
	)

	return dsn, func() { _ = pg.Terminate(ctx) }
}

func waitForPostgres(t *testing.T, dsn string) *sql.DB {
	t.Helper()

	dbx, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(15 * time.Second)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		err = dbx.PingContext(ctx)
		cancel()

		if err == nil {
			return dbx
		}

		if time.Now().After(deadline) {
			t.Fatalf("postgres not ready: %v", err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func newServer(t *testing.T, dsn string) (*httptest.Server, *sql.DB) {
	t.Helper()

	dbx := waitForPostgres(t, dsn)

	db.MustApplyMigrations(dbx)

	r := gin.Default()
	svc := service.Service{Notes: repo.NoteRepo{DB: dbx}}
	httpapi.Router{Svc: &svc}.Register(r)

	return httptest.NewServer(r), dbx
}

func truncateNotes(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`TRUNCATE notes RESTART IDENTITY`)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_CreateAndGetNote(t *testing.T) {
	dsn, stop := withPostgres(t)
	defer stop()

	srv, dbx := newServer(t, dsn)
	defer srv.Close()
	defer truncateNotes(t, dbx)

	resp, err := http.Post(
		srv.URL+"/notes",
		"application/json",
		strings.NewReader(`{"title":"Hello","content":"World"}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status=%d want=201", resp.StatusCode)
	}

	var created map[string]any
	b, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	_ = json.Unmarshal(b, &created)

	id := int64(created["id"].(float64))

	resp2, err := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, id))
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=200", resp2.StatusCode)
	}
}

func Test_GetNote_NotFound(t *testing.T) {
	dsn, stop := withPostgres(t)
	defer stop()

	srv, dbx := newServer(t, dsn)
	defer srv.Close()
	defer truncateNotes(t, dbx)

	resp, err := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, 999))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status=%d want=404", resp.StatusCode)
	}
}

func Test_DeleteNote(t *testing.T) {
	dsn, stop := withPostgres(t)
	defer stop()

	srv, dbx := newServer(t, dsn)
	defer srv.Close()
	defer truncateNotes(t, dbx)

	resp, _ := http.Post(
		srv.URL+"/notes",
		"application/json",
		strings.NewReader(`{"title":"t","content":"c"}`),
	)

	var created map[string]any
	b, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	_ = json.Unmarshal(b, &created)
	id := int64(created["id"].(float64))

	req, _ := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/notes/%d", srv.URL, id),
		nil,
	)
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=200", resp2.StatusCode)
	}

	resp3, _ := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, id))
	if resp3.StatusCode != http.StatusNotFound {
		t.Fatalf("status=%d want=404", resp3.StatusCode)
	}
}

func Test_ListNotes(t *testing.T) {
	dsn, stop := withPostgres(t)
	defer stop()

	srv, dbx := newServer(t, dsn)
	defer srv.Close()
	defer truncateNotes(t, dbx)

	for i := 1; i <= 2; i++ {
		_, err := http.Post(
			srv.URL+"/notes",
			"application/json",
			strings.NewReader(
				fmt.Sprintf(`{"title":"t%d","content":"c%d"}`, i, i),
			),
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	resp, err := http.Get(srv.URL + "/notes")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=200", resp.StatusCode)
	}

	var list []map[string]any
	b, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	_ = json.Unmarshal(b, &list)

	if len(list) != 2 {
		t.Fatalf("len=%d want=2", len(list))
	}
}
