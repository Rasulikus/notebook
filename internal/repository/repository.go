package repository

import (
	"context"

	"github.com/Rasulikus/notebook/internal/model"
)

type NoteRepository interface {
	Create(ctx context.Context, note *model.Note) error
	List(ctx context.Context) ([]model.Note, error)
}
