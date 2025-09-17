package tag

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

func (r *repo) Create(ctx context.Context, tag *model.Tag) error {
	err := r.db.NewInsert().Model(tag).Returning("*").Scan(ctx)
	return err
}

func (r *repo) List(ctx context.Context, userID int64) ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.db.NewSelect().Model(&tags).Where("user_id = ? OR user_id is NULL", userID).Scan(ctx); err != nil {
		return nil, err
	}
	return tags, nil
}
