package main

import (
	"log"
	"net/http"
	"os"

	httpapi "example.com/tech-ip-sem2/services/auth/internal/http"
	"example.com/tech-ip-sem2/services/auth/internal/service"
)

func main() {
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
	}

	authService := service.NewAuthService()
	handlers := httpapi.NewHandlers(authService)
	router := httpapi.NewRouter(handlers)

	log.Printf("Auth service started on :%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}