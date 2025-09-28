package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepo   repository.UserRepository
	jwtService *JWTService
}

func NewService(userRepo repository.UserRepository, jwtService *JWTService) *service {
	return &service{userRepo: userRepo, jwtService: jwtService}
}

func (s *service) Register(ctx context.Context, email, password, name string) error {
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

// должен вернуть access, refresh, userID и ошибку
func (s *service) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, userID int64, err error) {
	lowEmail := strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.GetByEmail(ctx, lowEmail)
	if err != nil {
		return "", "", 0, err
	}
	if user == nil {
		return "", "", 0, model.ErrBadRequest
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", 0, model.ErrBadRequest
	}

	accessToken, err = s.jwtService.CreateAccessToken(user.ID)
	if err != nil {
		return "", "", 0, model.ErrBadRequest
	}
	refreshToken, err = s.jwtService.CreateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", 0, model.ErrBadRequest
	}
	return accessToken, refreshToken, user.ID, nil
}

func (m *JWTService) Logout(ctx context.Context, refreshToken string) error {
	refreshTokenHash := generateRefreshTokenHash(refreshToken)
	err := m.sessionRepo.SetRevokedAtNow(ctx, refreshTokenHash)
	if err != nil {
		return model.ErrBadRequest
	}
	return nil
}
