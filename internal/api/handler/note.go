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

type NoteHandler struct {
	s service.NoteService
}

func NewNoteHandler(s service.NoteService) *NoteHandler {
	return &NoteHandler{s: s}
}

type NoteResp struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags"`
	UserID    int64     `json:"user_id"`
}

func toNoteResp(n *model.Note) NoteResp {
	tags := make([]string, 0, len(n.Tags))
	for _, tag := range n.Tags {
		tags = append(tags, tag.Name)
	}
	return NoteResp{
		ID:        n.ID,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Title:     n.Title,
		Text:      n.Text,
		Tags:      tags,
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

type CreateNoteReq struct {
	Title string  `json:"title" binding:"required,min=3,max=100"`
	Text  string  `json:"text"`
	Tags  []int64 `json:"tags"`
}

func (h *NoteHandler) Create(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

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
	newNote, err := h.s.Create(ctx, &n, req.Tags)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	c.JSON(http.StatusCreated, toNoteResp(newNote))
}

type NoteListQuery struct {
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
	Order  string `form:"order"`
}

func (h *NoteHandler) List(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	var q NoteListQuery

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
	userID := middleware.CurrentUserID(c)

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
		return
	}
	c.JSON(http.StatusOK, toNoteResp(note))
}

type UpdateByIDNoteReq struct {
	Title   *string  `json:"title" binding:"omitempty,min=1,max=100"`
	Text    *string  `json:"text"  binding:"omitempty,max=20000"`
	TagsIDs *[]int64 `json:"tags"`
}

// PATCH /notes/{id}
func (h *NoteHandler) UpdateByID(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

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
		status, pub := model.ToHTTP(err)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	ctx := c.Request.Context()
	r := &service.UpdateByIDNoteReq{
		Title:   req.Title,
		Text:    req.Text,
		TagsIDs: req.TagsIDs,
	}
	note, err := h.s.UpdateByID(ctx, userID, id, r)
	if err != nil {
		status, pub := model.ToHTTP(err)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	c.JSON(http.StatusOK, toNoteResp(note))
}

func (h *NoteHandler) DeleteByID(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

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
		return
	}

	c.JSON(http.StatusOK, "note is delete")
}
