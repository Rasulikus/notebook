package repository

import (
	"context"
	"time"

	"github.com/Rasulikus/notebook/internal/model"
)

type NoteRepository interface {
	Create(ctx context.Context, note *model.Note, tags []*model.Tag) (*model.Note, error)
	List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Note, error)
	GetByID(ctx context.Context, userID, id int64) (*model.Note, error)
	UpdateByID(ctx context.Context, userID, id int64, title, text *string, tagsIDs *[]int64) (*model.Note, error)
	DeleteTags(ctx context.Context, userID, id int64) error
	DeleteByID(ctx context.Context, userID, id int64) error
}

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	RotateRefreshToken(ctx context.Context, oldhash, newHash []byte, newExpiresAt time.Time) (*model.Session, error)
	SetRevokedAtNow(ctx context.Context, refreshTokenHash []byte) error
}

type TagRepository interface {
	Create(ctx context.Context, tag *model.Tag) error
	CreateTags(ctx context.Context, tags []*model.Tag) ([]*model.Tag, error)
	List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Tag, error)
	GetByID(ctx context.Context, userID, id int64) (*model.Tag, error)
	GetByIDs(ctx context.Context, userID int64, ids []int64) ([]*model.Tag, error)
	UpdateByID(ctx context.Context, userID int64, tag *model.Tag) (*model.Tag, error)
	DeleteByID(ctx context.Context, userID, id int64) error
}
