package service

import (
	"context"

	"github.com/Rasulikus/notebook/internal/model"
)

type NoteService interface {
	Create(ctx context.Context, note *model.Note) error
	List(ctx context.Context) ([]model.Note, error)
	GetByID(ctx context.Context, id int64) (*model.Note, error)
}

type AuthService interface {
	Register(ctx context.Context, email, password, name string) error
	Login(ctx context.Context, email, password string) error
}
