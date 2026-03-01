package main

import (
	"log"
	"net"
	"net/http"
	"os"

	httpapi "example.com/tech-ip-grpc/services/auth/internal/http"
	authgrpc "example.com/tech-ip-grpc/services/auth/internal/grpc"
	"example.com/tech-ip-grpc/services/auth/internal/service"
	"example.com/tech-ip-grpc/proto/authpb"

	"google.golang.org/grpc"
)

func main() {
	httpPort := os.Getenv("AUTH_PORT")
	if httpPort == "" {
		httpPort = "8081"
	}

	grpcPort := os.Getenv("AUTH_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	authService := service.NewAuthService()

	// --- HTTP server ---
	handlers := httpapi.NewHandlers(authService)
	router := httpapi.NewRouter(handlers)

	go func() {
		log.Printf("Auth HTTP service started on :%s", httpPort)
		if err := http.ListenAndServe(":"+httpPort, router); err != nil {
			log.Fatal(err)
		}
	}()

	// --- gRPC server ---
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(
		grpcServer,
		authgrpc.NewServer(authService),
	)

	log.Printf("Auth gRPC service started on :%s", grpcPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}