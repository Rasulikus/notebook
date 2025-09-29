package session

import (
	"context"
	"time"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
	"github.com/uptrace/bun"
)

type repo struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *repo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, session *model.Session) error {
	err := r.db.NewInsert().Model(session).Scan(ctx)
	return err
}

func (r *repo) RotateRefreshToken(ctx context.Context, oldHash, newHash []byte, newExpiresAt time.Time) (*model.Session, error) {
	sess := new(model.Session)

	_, err := r.db.NewUpdate().
		Model(sess).
		Set("refresh_token_hash = ?", newHash).
		Set("expires_at = ?", newExpiresAt).
		Where("refresh_token_hash = ?", oldHash).
		Where("revoked_at IS NULL").
		Where("expires_at > now()").
		Returning("*").
		Exec(ctx)
	return sess, repository.IsNoRowsError(err)
}

func (r *repo) SetRevokedAtNow(ctx context.Context, refreshTokenHash []byte) error {
	now := time.Now().UTC()

	res, err := r.db.NewUpdate().
		Model((*model.Session)(nil)).
		Set("revoked_at = ?", now).
		Where("refresh_token_hash = ?", refreshTokenHash).
		Where("revoked_at IS NULL").
		Where("expires_at > now()").
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
