// Package handler - Gin HTTP handlers for Tag CRUD.
package handler

import (
	"net/http"
	"strconv"

	"github.com/Rasulikus/notebook/internal/api/middleware"
	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/service"
	"github.com/gin-gonic/gin"
)

// TagHandler wires HTTP to TagService.
type TagHandler struct {
	s service.TagService
}

// NewTagHandler - constructor.
func NewTagHandler(s service.TagService) *TagHandler { return &TagHandler{s: s} }

// TagResp - public shape returned by the API.
type TagResp struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	UserID int64  `json:"user_id"`
}

// toTagResp - maps domain tag to API response.
func toTagResp(tag *model.Tag) TagResp {
	return TagResp{
		ID:     tag.ID,
		Name:   tag.Name,
		UserID: tag.UserID,
	}
}

// toTagsResp - maps slice of domain tags to []TagResp.
func toTagsResp(tags []model.Tag) []TagResp {
	out := make([]TagResp, 0, len(tags))
	for _, tag := range tags {
		out = append(out, toTagResp(&tag))
	}
	return out
}

// CreateTagReq request body for creating a tag.
type CreateTagReq struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}

// Create (POST /tags) creates a tag for current user; 201 + TagResp.
func (h *TagHandler) Create(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	var req CreateTagReq
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
	tag := model.Tag{Name: req.Name, UserID: userID}
	if err := h.s.Create(ctx, &tag); err != nil {
		status, pub := model.ToHTTP(err)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	c.JSON(http.StatusCreated, toTagResp(&tag))
}

// TagListQuery query params for listing tags.
// All query params are optional; if omitted, defaults are used by the service.
type TagListQuery struct {
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
	Order  string `form:"order"`
}

// List (GET /tags) returns user's tags; 200 + []TagResp.
func (h *TagHandler) List(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	var q TagListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	ctx := c.Request.Context()
	tags, err := h.s.List(ctx, userID, q.Limit, q.Offset, q.Order)
	if err != nil {
		status, pub := model.ToHTTP(err)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	c.JSON(http.StatusOK, toTagsResp(tags))
}

// GetByID (GET /tags/:id) loads one tag by id for current user; 200 + TagResp.
func (h *TagHandler) GetByID(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	ctx := c.Request.Context()
	tag, err := h.s.GetByID(ctx, userID, id)
	if err != nil {
		status, pub := model.ToHTTP(err)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	c.JSON(http.StatusOK, toTagResp(tag))
}

// UpdateByIDTagReq request body for updating tag name.
type UpdateByIDTagReq struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}

// UpdateByID (PATCH /tags/:id) updates tag name; 200 + TagResp.
func (h *TagHandler) UpdateByID(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	var req UpdateByIDTagReq
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
	tag := &model.Tag{ID: id, UserID: userID, Name: req.Name}
	updTag, err := h.s.UpdateByID(ctx, userID, tag)
	if err != nil {
		status, pub := model.ToHTTP(err)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	c.JSON(http.StatusOK, toTagResp(updTag))
}

// DeleteByID (DELETE /tags/:id) deletes a tag; 204 No Content.
func (h *TagHandler) DeleteByID(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	ctx := c.Request.Context()
	if err := h.s.DeleteByID(ctx, userID, id); err != nil {
		status, pub := model.ToHTTP(err)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	c.Status(http.StatusNoContent)
}
