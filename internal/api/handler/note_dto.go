package handler

import (
	"time"

	"github.com/Rasulikus/notebook/internal/model"
)

type CreateNoteReq struct {
	Title string `json:"title" binding:"required,min=1"`
	Text  string `json:"text"`
}

type NoteResp struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	UserID    int64     `json:"user_id"`
}

func toNoteResp(n model.Note) NoteResp {
	return NoteResp{
		ID:        n.ID,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Title:     n.Title,
		Text:      n.Text,
		UserID:    n.UserID,
	}
}

func toNotesResp(ns []model.Note) []NoteResp {
	out := make([]NoteResp, 0, len(ns))

	for _, n := range ns {
		out = append(out, toNoteResp(n))
	}
	return out
}
