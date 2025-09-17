package note

import (
	"context"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
)

type service struct {
	noteRepo repository.NoteRepository
}

func NewService(noteRepo repository.NoteRepository) *service {
	return &service{noteRepo: noteRepo}
}

func (s *service) Create(ctx context.Context, n *model.Note) error {
	return s.noteRepo.Create(ctx, n)
}

func (s *service) List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Note, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	if order == "" {
		order = "created_at"
	}
	return s.noteRepo.List(ctx, userID, limit, offset, order)
}

func (s *service) GetByID(ctx context.Context, userID, id int64) (*model.Note, error) {
	return s.noteRepo.GetByID(ctx, userID, id)
}

func (s *service) UpdateByID(ctx context.Context, userID int64, note *model.Note) error {
	return s.noteRepo.UpdateByID(ctx, userID, note)
}

func (s *service) DeleteByID(ctx context.Context, userID, id int64) error {
	return s.noteRepo.DeleteByID(ctx, userID, id)
}
