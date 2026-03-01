package authgrpc

import (
	"context"
	"errors"
	"log"
	"time"

	"example.com/tech-ip-grpc/proto/authpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnauthorized   = errors.New("unauthorized")
	ErrAuthUnavailable = errors.New("auth service unavailable")
)

type Client struct {
	client authpb.AuthServiceClient
	timeout time.Duration
}

func New(conn *grpc.ClientConn, timeout time.Duration) *Client {
	return &Client{
		client:  authpb.NewAuthServiceClient(conn),
		timeout: timeout,
	}
}

func (c *Client) Verify(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	log.Printf("calling grpc verify")
	
	_, err := c.client.Verify(ctx, &authpb.VerifyRequest{
		Token: token,
	})
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return ErrAuthUnavailable
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return ErrUnauthorized
	case codes.DeadlineExceeded:
		return ErrAuthUnavailable
	default:
		return ErrAuthUnavailable
	}
}