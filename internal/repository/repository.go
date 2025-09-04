package repository

import (
	"context"
	"time"

	"github.com/Rasulikus/notebook/internal/model"
)

type NoteRepository interface {
	Create(ctx context.Context, note *model.Note) error
	List(ctx context.Context) ([]model.Note, error)
	GetByID(ctx context.Context, id int64) (*model.Note, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	RotateRefreshTokenTx(ctx context.Context, oldhash, newHash []byte, newExpiresAt time.Time) (*model.Session, error)
	DeleteByUserID(ctx context.Context, userID int64) error
}
