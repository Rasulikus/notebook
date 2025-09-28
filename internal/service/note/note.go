package note

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
	"github.com/Rasulikus/notebook/internal/service"
)

type s struct {
	noteRepo repository.NoteRepository
	tagRepo  repository.TagRepository
}

func NewService(noteRepo repository.NoteRepository, tagRepo repository.TagRepository) *s {
	return &s{noteRepo: noteRepo, tagRepo: tagRepo}
}

func (s *s) Create(ctx context.Context, n *model.Note, tagsIDs []int64) (*model.Note, error) {
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

func (s *s) List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Note, error) {
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

func (s *s) GetByID(ctx context.Context, userID, id int64) (*model.Note, error) {
	return s.noteRepo.GetByID(ctx, userID, id)
}

func (s *s) UpdateByID(ctx context.Context, userID, id int64, req *service.UpdateByIDNoteReq) (*model.Note, error) {
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

func (s *s) DeleteByID(ctx context.Context, userID, id int64) error {
	return s.noteRepo.DeleteByID(ctx, userID, id)
}
