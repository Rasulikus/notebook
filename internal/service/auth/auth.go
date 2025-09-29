package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type TokenConfig struct {
	Secret      []byte
	AccessTTL   time.Duration // lifetime for access JWT
	RefreshTTL  time.Duration // lifetime for refresh token/session
	SessionRepo repository.SessionRepository
}

type Service struct {
	userRepo     repository.UserRepository
	tokenManager *tokenManager
}

func NewService(userRepo repository.UserRepository, cfg TokenConfig) *Service {
	return &Service{userRepo: userRepo, tokenManager: newTokenManager(cfg)}
}

// Register creates a new user.
func (s *Service) Register(ctx context.Context, email, password, name string) error {
	lowEmail := strings.ToLower(strings.TrimSpace(email))
	_, err := s.userRepo.GetByEmail(ctx, lowEmail)
	if err == nil {
		return model.ErrConflict

	}
	if !errors.Is(err, model.ErrNotFound) {
		return err
	}

	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}

	newUser := &model.User{
		Email:        lowEmail,
		PasswordHash: string(hash),
		Name:         name,
	}
	err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return err
	}
	return nil
}

// Login authenticates user and issues tokens.
func (s *Service) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, userID int64, err error) {
	lowEmail := strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.GetByEmail(ctx, lowEmail)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", "", 0, model.ErrWrongCredentials
		}
		return "", "", 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", 0, model.ErrWrongCredentials
	}

	accessToken, err = s.tokenManager.CreateAccessToken(user.ID)
	if err != nil {
		return "", "", 0, err
	}
	refreshToken, err = s.tokenManager.CreateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", 0, err
	}
	return accessToken, refreshToken, user.ID, nil
}

// Refresh rotates refresh token and returns a new access + refresh.
func (s *Service) Refresh(ctx context.Context, oldRefreshToken string) (string, string, error) {
	return s.tokenManager.RotateRefreshToken(ctx, oldRefreshToken)
}

func (s *Service) ParseAccessToken(token string) (int64, error) {
	return s.tokenManager.ParseAccessToken(token)
}

// Logout revokes refresh token by setting revoked_at now.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	refreshTokenHash := s.tokenManager.generateRefreshTokenHash(refreshToken)
	return s.tokenManager.SessionRepo.SetRevokedAtNow(ctx, refreshTokenHash)
}
