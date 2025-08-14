package note

import (
	"context"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/uptrace/bun"
)

type repo struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *repo {
	db.RegisterModel((*model.NoteTag)(nil))
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, n *model.Note) error {
	err := r.db.NewInsert().Model(n).Scan(ctx, n)
	return err
}

func (r *repo) List(ctx context.Context) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.NewSelect().Model(&notes).Scan(ctx)
	return notes, err
}
