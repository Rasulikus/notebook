package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Note struct {
	bun.BaseModel `bun:"table:notes"`
	ID            int64     `json:"id" bun:"id,pk,autoincrement"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,notnull,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,notnull,nullzero,default:current_timestamp"`
	Title         string    `json:"title" bun:"title,notnull"`
	Text          string    `json:"text" bun:"text,"`

	Tags []*Tag `bun:"m2m:notes_tags,join:Note=Tag"`

	UserID int64 `json:"user_id" bun:"user_id,notnull"`
}
