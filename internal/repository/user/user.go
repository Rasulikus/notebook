package user

import (
	"context"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/uptrace/bun"
)

type repo struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *repo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, user *model.User) error {
	err := r.db.NewInsert().Model(user).Scan(ctx, user)
	return err
}

func (r *repo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().Model(user).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
