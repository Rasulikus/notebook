package auth

import (
	"context"
	"database/sql"
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
	existedUser, err := s.userRepo.GetByEmail(ctx, lowEmail)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if existedUser != nil {
		return model.ErrUserAlreadyExists
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
func (s *service) Login(ctx context.Context, email, password string) (string, string, int64, error) {
	lowEmail := strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.GetByEmail(ctx, lowEmail)
	if err != nil {
		return "", "", 0, err
	}
	if user == nil {
		return "", "", 0, model.ErrWrongCredetials
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", 0, model.ErrWrongCredetials
	}

	access, err := s.jwtService.CreateAccessToken(user.ID)
	if err != nil {
		return "", "", 0, model.ErrWrongCredetials
	}
	refresh, err := s.jwtService.CreateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", 0, model.ErrWrongCredetials
	}
	return access, refresh, user.ID, nil
}

func (m *JWTService) Logout(ctx context.Context, refreshToken string) error {
	refreshTokenHash := generateRefreshTokenHash(refreshToken)
	err := m.sessionRepo.SetRevokedAtNow(ctx, refreshTokenHash)
	if err != nil {
		return model.ErrInvalidToken
	}
	return nil
}
