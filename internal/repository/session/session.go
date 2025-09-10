package session

import (
	"context"
	"time"

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
	err := r.db.NewInsert().Model(session).Scan(ctx)
	return err
}

func (r *repo) RotateRefreshTokenTx(ctx context.Context, oldhash, newHash []byte, newExpiresAt time.Time) (*model.Session, error) {
	session := new(model.Session)

	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

		err := tx.NewSelect().
			Model(session).
			Where("refresh_token_hash = ?", oldhash).
			Where("revoked_at IS NULL").
			Where("expires_at > now()").
			For("UPDATE").
			Scan(ctx)
		if err != nil {
			return err
		}

		res, err := tx.NewUpdate().
			Model(session).
			Set("refresh_token_hash = ?", newHash).
			WherePK().
			Returning("*").
			Exec(ctx)
		if err != nil {
			return err
		}
		aff, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if aff != 1 {
			return model.ErrInvalidToken
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return session, nil
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
	aff, _ := res.RowsAffected()
	if aff != 1 {
		return model.ErrInvalidToken
	}
	return nil
}
