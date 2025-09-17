package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Rasulikus/notebook/internal/api/middleware"
	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/service"
	"github.com/gin-gonic/gin"
)

type CreateNoteReq struct {
	Title string `json:"title" binding:"required,min=3,max=100"`
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

type ListQuery struct {
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
	Order  string `form:"order"`
}

type UpdateByIDNoteReq struct {
	Title string `json:"title" binding:"required,min=3,max=100"`
	Text  string `json:"text"`
}

func toNoteResp(n *model.Note) NoteResp {
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
		out = append(out, toNoteResp(&n))
	}
	return out
}

type NoteHandler struct {
	s service.NoteService
}

func NewNoteHandler(s service.NoteService) *NoteHandler {
	return &NoteHandler{s: s}
}

func (h *NoteHandler) Create(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		status, pub := model.ToHTTP(model.ErrUnauthorized)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	var req CreateNoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if vErr, as := model.AsValidationError(req, err); as {
			status, pub := model.ToHTTP(vErr)
			c.AbortWithStatusJSON(status, pub)
			return
		}
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	ctx := c.Request.Context()
	n := model.Note{Title: req.Title, Text: req.Text, UserID: userID}
	if err := h.s.Create(ctx, &n); err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
	}
	c.JSON(http.StatusCreated, toNoteResp(&n))
}

func (h *NoteHandler) List(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		status, pub := model.ToHTTP(model.ErrUnauthorized)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	var q ListQuery

	if err := c.ShouldBindQuery(&q); err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	ctx := c.Request.Context()
	notes, err := h.s.List(ctx, userID, q.Limit, q.Offset, q.Order)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	resp := toNotesResp(notes)
	c.JSON(http.StatusOK, resp)
}

func (h *NoteHandler) GetByID(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		status, pub := model.ToHTTP(model.ErrUnauthorized)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	ctx := c.Request.Context()
	note, err := h.s.GetByID(ctx, userID, id)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
	}
	c.JSON(http.StatusOK, toNoteResp(note))
}

// PATCH /notes/{id}
func (h *NoteHandler) UpdateByID(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		status, pub := model.ToHTTP(model.ErrUnauthorized)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	var req UpdateByIDNoteReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		if vErr, as := model.AsValidationError(req, err); as {
			status, pub := model.ToHTTP(vErr)
			c.AbortWithStatusJSON(status, pub)
			return
		}
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	ctx := c.Request.Context()
	n := model.Note{
		Title:  req.Title,
		Text:   req.Text,
		ID:     id,
		UserID: userID,
	}
	h.s.UpdateByID(ctx, userID, &n)
	c.JSON(http.StatusOK, toNoteResp(&n))
}

func (h *NoteHandler) DeleteByID(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		status, pub := model.ToHTTP(model.ErrUnauthorized)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	ctx := c.Request.Context()
	err = h.s.DeleteByID(ctx, userID, id)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
	}
	c.JSON(http.StatusOK, "note is delete")
}
