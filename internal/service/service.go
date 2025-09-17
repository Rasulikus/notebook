package service

import (
	"context"

	"github.com/Rasulikus/notebook/internal/model"
)

type NoteService interface {
	Create(ctx context.Context, note *model.Note) error
	List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Note, error)
	GetByID(ctx context.Context, userID, id int64) (*model.Note, error)
	UpdateByID(ctx context.Context, userID int64, note *model.Note) error
	DeleteByID(ctx context.Context, userID, id int64) error
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
