package handler

import (
	"net/http"
	"time"

	"github.com/Rasulikus/notebook/internal/api/middleware"
	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/service"
	"github.com/gin-gonic/gin"
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

type NoteHandler struct {
	s service.NoteService
}

func NewNoteHanlder(s service.NoteService) *NoteHandler {
	return &NoteHandler{s: s}
}

func (h *NoteHandler) Create(c *gin.Context) {
	var req CreateNoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
	}
	ctx := c.Request.Context()
	n := model.Note{Title: req.Title, Text: req.Text, UserID: userID}
	if err := h.s.Create(ctx, &n); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, NoteResp{
		ID:        n.ID,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Title:     n.Title,
		Text:      n.Text,
		UserID:    n.UserID,
	})
}

func (h *NoteHandler) List(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
	}
	ctx := c.Request.Context()
	notes, err := h.s.List(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := toNotesResp(notes)
	c.JSON(http.StatusOK, resp)
}
