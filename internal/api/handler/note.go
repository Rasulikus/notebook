package handler

import (
	"net/http"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/service"
	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	s service.NoteService
}

func NewNoteHanlder(s service.NoteService) *NoteHandler {
	return &NoteHandler{s: s}
}

func (h *NoteHandler) RegisterNotes(r *gin.RouterGroup) {
	r.POST("/", h.create)
	r.GET("/", h.list)
}

func (h *NoteHandler) create(c *gin.Context) {
	var req CreateNoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()

	n := model.Note{Title: req.Title, Text: req.Text}
	if err := h.s.Create(ctx, &n); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, req)
}

func (h *NoteHandler) list(c *gin.Context) {
	ctx := c.Request.Context()
	notes, err := h.s.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := toNotesResp(notes)
	c.JSON(http.StatusOK, resp)
}
