package main

import (
	"log"
	"net/http"
	"os"
	"time"

	httpapi "example.com/tech-ip-grpc/services/tasks/internal/http"
	"example.com/tech-ip-grpc/services/tasks/internal/service"
	"example.com/tech-ip-grpc/services/tasks/internal/client/authgrpc"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}

	authAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authAddr == "" {
		log.Fatal("AUTH_GRPC_ADDR is required")
	}

	// --- gRPC connection ---
	conn, err := grpc.Dial(
		authAddr,
		grpc.WithInsecure(), // допустимо для учебного ПЗ
	)
	if err != nil {
		log.Fatalf("failed to connect to auth grpc: %v", err)
	}
	defer conn.Close()

	authClient := authgrpc.New(conn, 2*time.Second)

	// --- HTTP server ---
	taskService := service.NewTaskService()
	handlers := httpapi.NewHandlers(taskService)
	router := httpapi.NewRouter(handlers, authClient)

	log.Printf("Tasks service started on :%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}