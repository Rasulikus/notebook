package model

import "github.com/uptrace/bun"

type NoteTag struct {
	bun.BaseModel `bun:"table:notes_tags"`
	NoteID        int64 `json:"note_id" bun:",pk"`
	TagID         int64 `json:"tag_id" bun:",pk"`
}
