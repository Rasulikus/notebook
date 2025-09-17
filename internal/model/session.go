package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel    `bun:"table:sessions"`
	ID               int64     `bun:"id,pk,autoincrement"`
	UserID           int64     `bun:"user_id,notnull,nullzero"`
	RefreshTokenHash []byte    `bun:"refresh_token_hash,type:bytea,unique,notnull" json:"-"`
	ExpiresAt        time.Time `bun:"expires_at,notnull,nullzero"`
	CreatedAt        time.Time `bun:"created_at,notnull,nullzero,default:current_timestamp"`
	RevokedAt        time.Time `bun:"revoked_at,nullzero"`
}
