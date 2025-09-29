package service

import (
	"context"

	"github.com/Rasulikus/notebook/internal/model"
)

type UpdateByIDNoteReq struct {
	Title   *string
	Text    *string
	TagsIDs *[]int64
}

type NoteService interface {
	Create(ctx context.Context, n *model.Note, TagsIDs []int64) (*model.Note, error)
	List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Note, error)
	GetByID(ctx context.Context, userID, id int64) (*model.Note, error)
	UpdateByID(ctx context.Context, userID, id int64, req *UpdateByIDNoteReq) (*model.Note, error)
	DeleteByID(ctx context.Context, userID, id int64) error
}

type AuthService interface {
	Register(ctx context.Context, email, password, name string) error
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, userID int64, err error)
	Refresh(ctx context.Context, oldRefreshToken string) (string, string, error)
	ParseAccessToken(token string) (int64, error)
	Logout(ctx context.Context, refreshToken string) error
}

type JWTService interface {
	CreateAccessToken(userID int64) (string, error)
	CreateRefreshToken(ctx context.Context, userID int64) (string, error)
	RotateRefreshToken(ctx context.Context, oldRefresh string) (string, string, error)
	ParseAccessToken(token string) (int64, error)
}

type TagService interface {
	Create(ctx context.Context, tag *model.Tag) error
	List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Tag, error)
	GetByID(ctx context.Context, userID, id int64) (*model.Tag, error)
	UpdateByID(ctx context.Context, userID int64, tag *model.Tag) (*model.Tag, error)
	DeleteByID(ctx context.Context, userID, id int64) error
}
