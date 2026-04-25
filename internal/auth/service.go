package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, username, password string, role UserRole) error {
	user := &User{
		Username: username,
		Password: password,
		Role:     role,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, username, password string) (string, time.Time, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", time.Time{}, err
	}

	if password != user.Password {
		return "", time.Time{}, errors.New("Incorrect password")
	}

	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	fmt.Println(user.ID)

	session := &Session{
		Token:     token,
		UserId:    user.ID,
		ExpiresAt: expiresAt,
	}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	return s.repo.DeleteSession(ctx, token)
}

func (s *Service) Authenticate(ctx context.Context, token string) (*User, error) {
	return s.repo.GetUserBySessionToken(ctx, token)
}

// func (s *Service) HasRole(user *User, requiredRole string) bool
