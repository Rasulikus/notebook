package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int64     `json:"id" bun:"id,pk,autoincrement"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt     time.Time `json:"deleted_at" bun:"deleted_at,soft_delete"`
	Email         string    `json:"email" bun:"email,unique,notnull"`
	PasswordHash  string    `json:"-" bun:"password_hash,notnull"`
	Name          string    `json:"name" bun:"name,notnull"`
}
