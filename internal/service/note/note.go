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

func (s *service) List(ctx context.Context) ([]model.Note, error) {
	return s.noteRepo.List(ctx)
}
