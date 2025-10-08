package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/pz3-http/internal/api"
	"example.com/pz3-http/internal/storage"
)

func main() {
	store := storage.NewMemoryStore()
	h := api.NewHandlers(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		api.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Коллекция
	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	// Элемент
	mux.HandleFunc("GET /tasks/", h.GetTask)
	mux.HandleFunc("PATCH /tasks/", h.MarkTaskDone)
	mux.HandleFunc("DELETE /tasks/", h.DeleteTask)

	// Подключаем логирование и CORS
	handler := api.Logging(api.CORS(mux))

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Канал для сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("listening on", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-stop // Ожидание сигнала

	log.Println("shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server stopped")
}
