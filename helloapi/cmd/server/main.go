package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Текстовый ответ
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user{
			ID:   uuid.NewString(), // теперь реальный UUID
			Name: "Gopher",
		})
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Status string `json:"status"`
			Time   string `json:"time"`
		}{
			Status: "ok",
			Time:   time.Now().Format(time.RFC3339),
		}
		_ = json.NewEncoder(w).Encode(response)
	})

	addr := ":" + port
	log.Printf("Starting on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
