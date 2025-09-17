package repository

import (
	"context"
	"time"

	"github.com/Rasulikus/notebook/internal/model"
)

type NoteRepository interface {
	Create(ctx context.Context, note *model.Note) error
	List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Note, error)
	GetByID(ctx context.Context, userID, id int64) (*model.Note, error)
	UpdateByID(ctx context.Context, userID int64, note *model.Note) error
	DeleteByID(ctx context.Context, userID, id int64) error
}

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	RotateRefreshTokenTx(ctx context.Context, oldhash, newHash []byte, newExpiresAt time.Time) (*model.Session, error)
	SetRevokedAtNow(ctx context.Context, refreshTokenHash []byte) error
}
