package tag

import "github.com/uptrace/bun"

type repo struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *repo {
	return &repo{db: db}
}
