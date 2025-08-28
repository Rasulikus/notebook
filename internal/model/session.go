package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel    `bun:"table:sessions"`
	ID               int64     `bun:",pk,autoincrement"`
	UserID           int64     `bun:",notnull"`
	RefreshTokenHash []byte    `bun:"refresh_token_hash,type:bytea,unique,notnull" json:"-"`
	ExpiresAt        time.Time `bun:",notnull"`
	CreatedAt        time.Time `bun:",nullzero,default:current_timestamp"`
}
