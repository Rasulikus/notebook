package model

import "github.com/uptrace/bun"

type NoteTag struct {
	bun.BaseModel `bun:"table:notes_tags"`
	NoteID        int64 `json:"note_id" bun:"note_id,pk"`
	Note          *Note `bun:"rel:belongs-to,join:note_id=id"`
	TagID         int64 `json:"tag_id" bun:"tag_id,pk"`
	Tag           *Tag  `bun:"rel:belongs-to,join:tag_id=id"`
}
