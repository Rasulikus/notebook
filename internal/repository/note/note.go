package note

import (
	"context"
	"errors"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/uptrace/bun"
)

var (
	ErrRowsAffected = errors.New("affected rows more than need necessary")
)

type repo struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *repo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, note *model.Note) error {
	_, err := r.db.NewInsert().Model(note).Returning("*").Exec(ctx)
	return err
}

func (r *repo) List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.NewSelect().
		Model(&notes).
		Where("user_id = ?", userID).
		Order(order).
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return notes, err
}

func (r *repo) GetByID(ctx context.Context, userID, id int64) (*model.Note, error) {
	note := new(model.Note)
	err := r.db.NewSelect().Model(note).Where("id = ? AND user_id = ?", id, userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (r *repo) UpdateByID(ctx context.Context, userID int64, note *model.Note) error {
	res, err := r.db.NewUpdate().
		Model(note).
		Set("title = ?", note.Title).
		Set("text = ?", note.Text).
		Set("updated_at = now()").
		Where("id = ? AND user_id = ?", note.ID, userID).
		Returning("*").
		Exec(ctx)
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

func (r *repo) DeleteByID(ctx context.Context, userID, id int64) error {
	_, err := r.db.NewDelete().Model((*model.Note)(nil)).Where("id = ? AND user_id = ?", id, userID).Exec(ctx)
	return err
}
