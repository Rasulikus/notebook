package service

import (
	"context"

	"github.com/Rasulikus/notebook/internal/model"
)

type NoteService interface {
	Create(ctx context.Context, note *model.Note) error
	List(ctx context.Context, userID int64) ([]model.Note, error)
	GetByID(ctx context.Context, id int64) (*model.Note, error)
}

type AuthService interface {
	Register(ctx context.Context, email, password, name string) error
	// должен вернуть access, refresh, userID и ошибку
	Login(ctx context.Context, email, password string) (string, string, int64, error)
}

type JWTService interface {
	CreateAccessToken(userID int64) (string, error)
	CreateRefreshToken(ctx context.Context, userID int64) (string, error)
	RotateRefreshToken(ctx context.Context, oldRefresh string) (string, string, error)
	ParseAccessToken(token string) (int64, error)
	Logout(ctx context.Context, refreshToken string) error
}
