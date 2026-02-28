package main

import (
	"log"
	"net/http"
	"os"
	"time"

	httpapi "example.com/tech-ip-sem2/services/tasks/internal/http"
	"example.com/tech-ip-sem2/services/tasks/internal/service"
	"example.com/tech-ip-sem2/services/tasks/internal/client/authclient"
	"example.com/tech-ip-sem2/shared/httpx"
)

func main() {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}

	authURL := os.Getenv("AUTH_BASE_URL")
	if authURL == "" {
		log.Fatal("AUTH_BASE_URL is required")
	}

	taskService := service.NewTaskService()
	handlers := httpapi.NewHandlers(taskService)

	httpClient := httpx.New(3 * time.Second)
	authClient := authclient.New(authURL, httpClient)

	router := httpapi.NewRouter(handlers, authClient)

	log.Printf("Tasks service started on :%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}