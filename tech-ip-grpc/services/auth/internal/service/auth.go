package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Login — упрощённая авторизация
func (s *AuthService) Login(username, password string) (string, error) {
	// Учебная логика
	if username == "student" && password == "student" {
		return "demo-token", nil
	}
	return "", ErrInvalidCredentials
}

// VerifyToken — проверка токена
func (s *AuthService) VerifyToken(token string) (string, error) {
	if token == "demo-token" {
		return "student", nil
	}
	return "", ErrInvalidToken
}