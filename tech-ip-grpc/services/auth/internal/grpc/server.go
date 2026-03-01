package grpc

import (
	"context"
	"log"

	"example.com/tech-ip-grpc/proto/authpb"
	"example.com/tech-ip-grpc/services/auth/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server реализует authpb.AuthServiceServer
type Server struct {
	authpb.UnimplementedAuthServiceServer
	auth *service.AuthService
}

func NewServer(auth *service.AuthService) *Server {
	return &Server{auth: auth}
}

// Verify — gRPC-метод проверки токена
func (s *Server) Verify(
	ctx context.Context,
	req *authpb.VerifyRequest,
) (*authpb.VerifyResponse, error) {

	log.Printf("grpc verify called, token=%s", req.Token)
	
	token := req.GetToken()
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is empty")
	}

	subject, err := s.auth.VerifyToken(token)
	if err != nil {
		log.Printf("grpc verify failed: %v", err)
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &authpb.VerifyResponse{
		Valid:   true,
		Subject: subject,
	}, nil
}