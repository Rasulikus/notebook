package session

import (
	"context"
	"time"

	"github.com/Rasulikus/notebook/internal/errs"
	"github.com/Rasulikus/notebook/internal/model"
	"github.com/uptrace/bun"
)

type repo struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *repo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, session *model.Session) error {
	err := r.db.NewInsert().Model(session).Scan(ctx, session)
	return err
}

func (r *repo) FindByHash(ctx context.Context, hash []byte) (*model.Session, error) {
	session := new(model.Session)
	err := r.db.NewSelect().Model(session).Where("refresh_token_hash = ?", hash).Scan(ctx, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *repo) Rotate(ctx context.Context, oldHash, newHash []byte, newExp time.Time) error {
	res, err := r.db.NewUpdate().
		Model((*model.Session)(nil)).
		Set("refresh_token_hash = ?", newHash).
		Set("expires_at = ?", newExp).
		Where("refresh_token_hash = ?", oldHash).
		Where("expires_at > NOW()").
		Exec(ctx)

	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff != 1 {
		return errs.ErrInvalidToken
	}
	return nil
}

func (r *repo) DeleteByUserID(ctx context.Context, userID int64) error {
	err := r.db.NewDelete().Model((*model.Session)(nil)).Where("user_id = ?", userID).Scan(ctx)
	return err
}
