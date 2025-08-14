package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int64     `json:"id" bun:",pk,autoincrement"`
	CreatedAt     time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:",nullzero"`
	DeletedAt     time.Time `json:"deleted_at" bun:",soft_delete,nullzero"`
	Email         string    `json:"email" bun:",notnull"`
	PasswordHash  string    `json:"-" bun:",notnull"`
	UserName      string    `json:"username" bun:"username,notnull"`
	Notes         []*Note   `json:"notes" bun:"rel:has-many,join:id=user_id"`
	Tags          []*Tag    `json:"tags" bun:"rel:has-many,join:id=user_id"`
}
