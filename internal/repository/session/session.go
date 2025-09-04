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

func (r *repo) DeleteByUserID(ctx context.Context, userID int64) error {
	_, err := r.db.NewDelete().Model((*model.Session)(nil)).Where("user_id = ?", userID).Exec(ctx)
	return err
}
