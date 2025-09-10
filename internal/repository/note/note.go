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
	db.RegisterModel((*model.NoteTag)(nil))
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, note *model.Note) error {
	_, err := r.db.NewInsert().Model(note).Exec(ctx)
	return err
}

func (r *repo) List(ctx context.Context, userID int64) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.NewSelect().Model(&notes).Where("user_id = ?", userID).Scan(ctx)
	return notes, err
}

func (r *repo) GetByID(ctx context.Context, id int64) (*model.Note, error) {
	note := new(model.Note)
	err := r.db.NewSelect().Model(note).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return note, nil
}
