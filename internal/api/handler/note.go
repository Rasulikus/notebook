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
	var req model.Note
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.s.Create(c, &req)
	c.JSON(http.StatusCreated, req)
}

func (h *NoteHandler) list(c *gin.Context) {
	notes, err := h.s.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notes)
}
