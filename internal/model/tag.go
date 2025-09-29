package model

import "github.com/uptrace/bun"

type Tag struct {
	bun.BaseModel `bun:"table:tags"`
	ID            int64  `json:"id" bun:"id,pk,autoincrement"`
	Name          string `json:"name" bun:"name,notnull"`

	UserID int64 `json:"user_id" bun:"user_id,nullzero"`
}
