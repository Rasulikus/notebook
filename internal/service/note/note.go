// Package note provides business logic for notes.
package note

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
	"github.com/Rasulikus/notebook/internal/service"
)

// Package note provides business logic for notes.
type Service struct {
	noteRepo repository.NoteRepository
	tagRepo  repository.TagRepository
}

// Package note provides business logic for notes.
func NewService(noteRepo repository.NoteRepository, tagRepo repository.TagRepository) *Service {
	return &Service{noteRepo: noteRepo, tagRepo: tagRepo}
}

// Create builds a note and attaches tags by IDs.
func (s *Service) Create(ctx context.Context, n *model.Note, tagsIDs []int64) (*model.Note, error) {
	for _, tagID := range tagsIDs {
		tag, err := s.tagRepo.GetByID(ctx, n.UserID, tagID)
		if err != nil {
			err = repository.IsNoRowsError(err)
			return nil, err
		}
		n.Tags = append(n.Tags, tag)
	}
	return s.noteRepo.Create(ctx, n, n.Tags)
}

// List returns user notes with sane defaults.
func (s *Service) List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Note, error) {
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

// GetByID returns a single note owned by the user.
func (s *Service) GetByID(ctx context.Context, userID, id int64) (*model.Note, error) {
	return s.noteRepo.GetByID(ctx, userID, id)
}

// UpdateByID applies partial changes and optionally replaces tags.
func (s *Service) UpdateByID(ctx context.Context, userID, id int64, req *service.UpdateByIDNoteReq) (*model.Note, error) {
	if req.TagsIDs != nil {
		_, err := s.tagRepo.GetByIDs(ctx, userID, *req.TagsIDs)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return nil, fmt.Errorf("tag not found: %w", model.ErrNotFound)
			}
			return nil, err
		}
	}

	note, err := s.noteRepo.UpdateByID(ctx, userID, id, req.Title, req.Text, req.TagsIDs)
	if err != nil {
		return nil, err
	}
	return note, nil
}

// DeleteByID removes a note owned by the user.
func (s *Service) DeleteByID(ctx context.Context, userID, id int64) error {
	return s.noteRepo.DeleteByID(ctx, userID, id)
}
