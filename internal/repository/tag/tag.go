package tag

import (
	"context"
	"fmt"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
	"github.com/uptrace/bun"
)

type Repo struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, tag *model.Tag) error {
	_, err := r.db.NewInsert().Model(tag).Returning("*").Exec(ctx)
	return err
}

func (r *Repo) CreateTags(ctx context.Context, tags []*model.Tag) ([]*model.Tag, error) {
	if len(tags) == 0 {
		return tags, nil

	}
	for _, t := range tags {
		if t == nil {
			return nil, fmt.Errorf("nil tag in slice")
		}
		if t.Name == "" {
			return nil, fmt.Errorf("tag name is empty")
		}
	}
	_, err := r.db.NewInsert().Model(&tags).On("CONFLICT (user_id, name) DO NOTHING").Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *Repo) List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.NewSelect().
		Model(&tags).
		Where("user_id = ? OR user_id is NULL ", userID).
		Order(order).
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *Repo) GetByID(ctx context.Context, userID, id int64) (*model.Tag, error) {
	tag := new(model.Tag)
	err := r.db.NewSelect().Model(tag).Where("id = ? AND user_id = ?", id, userID).Scan(ctx)
	if err != nil {
		return nil, repository.IsNoRowsError(err)
	}
	return tag, nil
}

func (r *Repo) GetByIDs(ctx context.Context, userID int64, ids []int64) ([]*model.Tag, error) {
	var tags []*model.Tag
	if len(ids) == 0 {
		return tags, nil
	}
	err := r.db.NewSelect().Model(&tags).Where("id IN (?) AND user_id = ?", bun.In(ids), userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	if len(tags) != len(ids) {
		return nil, model.ErrNotFound
	}
	return tags, nil
}

func (r *Repo) UpdateByID(ctx context.Context, userID int64, tag *model.Tag) (*model.Tag, error) {
	res, err := r.db.NewUpdate().
		Model(tag).
		Set("name = ?", tag.Name).
		Where("id = ? AND user_id = ?", tag.ID, userID).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if aff == 0 {
		return nil, model.ErrNotFound
	}
	return tag, nil
}

func (r *Repo) DeleteByID(ctx context.Context, userID, id int64) error {
	res, err := r.db.NewDelete().Model((*model.Tag)(nil)).Where("id = ? AND user_id = ?", id, userID).Exec(ctx)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return model.ErrNotFound
	}
	return nil
}
