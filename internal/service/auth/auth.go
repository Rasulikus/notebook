package auth

import (
	"context"
	"strings"

	"github.com/Rasulikus/notebook/internal/errs"
	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) *service {
	return &service{userRepo: userRepo}
}

func (s *service) Register(ctx context.Context, email, password, name string) error {
	lowEmail := strings.ToLower(strings.TrimSpace(email))
	existedUser, err := s.userRepo.FindByEmail(ctx, lowEmail)
	if err != nil {
		return err
	}
	if existedUser != nil {
		return errs.ErrUserAlreadyExists
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

func (s *service) Login(ctx context.Context, email, password string) error {
	lowEmail := strings.ToLower(strings.TrimSpace(email))
	existedUser, err := s.userRepo.FindByEmail(ctx, lowEmail)
	if err != nil {
		return err
	}
	if existedUser == nil {
		return errs.ErrWrongCredetials
	}
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.PasswordHash), []byte(password))
	if err != nil {
		return errs.ErrWrongCredetials
	}
	return nil
}
