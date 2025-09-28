package tag

import (
	"context"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
	"github.com/Rasulikus/notebook/internal/service"
)

var _ service.TagService = (*Service)(nil)

type Service struct {
	tagRepo repository.TagRepository
}

func NewService(tagRepo repository.TagRepository) *Service {
	return &Service{tagRepo: tagRepo}
}

func (s *Service) Create(ctx context.Context, tag *model.Tag) error {
	return s.tagRepo.Create(ctx, tag)
}

func (s *Service) List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Tag, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	if order == "" {
		order = "id"
	}
	return s.tagRepo.List(ctx, userID, limit, offset, order)
}

func (s *Service) GetByID(ctx context.Context, userID, id int64) (*model.Tag, error) {
	return s.tagRepo.GetByID(ctx, userID, id)
}

func (s *Service) UpdateByID(ctx context.Context, userID int64, tag *model.Tag) (*model.Tag, error) {
	tag, err := s.tagRepo.UpdateByID(ctx, userID, tag)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *Service) DeleteByID(ctx context.Context, userID, id int64) error {
	return s.tagRepo.DeleteByID(ctx, userID, id)
}
