package model

import "github.com/uptrace/bun"

type Tag struct {
	bun.BaseModel `bun:"table:tags"`
	ID            int64  `json:"id" bun:",pk,autoincrement"`
	Name          string `json:"name" bun:",notnull"`

	UserID int64 `json:"user_id" bun:",nullzero"`
	User   *User `json:"user" bun:"rel:belongs-to,join:user_id=id"`

	Notes []*Note `json:"notes" bun:"m2m:notes_tags,join:Tag=Note"`
}
