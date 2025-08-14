package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Note struct {
	bun.BaseModel `bun:"table:notes"`
	ID            int64     `json:"id" bun:",pk,autoincrement"`
	CreatedAt     time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:",nullzero"`
	Title         string    `json:"title" bun:",notnull"`
	Text          string    `json:"text"`

	Tags []*Tag `bun:"m2m:notes_tags,join:Note=Tag"`

	UserID int64 `json:"user_id"`
	User   *User `json:"user" bun:"rel:belongs-to,join:user_id=id"`
}
