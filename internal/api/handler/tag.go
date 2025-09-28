package handler

import (
	"net/http"
	"strconv"

	"github.com/Rasulikus/notebook/internal/api/middleware"
	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/service"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	s service.TagService
}

func NewTagHandler(s service.TagService) *TagHandler {
	return &TagHandler{s: s}
}

type TagResp struct {
	ID     int64
	Name   string
	UserID int64
}

func toTagResp(tag *model.Tag) TagResp {
	return TagResp{
		ID:     tag.ID,
		Name:   tag.Name,
		UserID: tag.UserID,
	}
}

func toTagsResp(tags []model.Tag) []TagResp {
	out := make([]TagResp, 0, len(tags))
	for _, tag := range tags {
		out = append(out, toTagResp(&tag))
	}
	return out
}

type CreateTagReq struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}

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
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	c.JSON(http.StatusCreated, toTagResp(&tag))
}

type TagListQuery struct {
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
	Order  string `form:"order"`
}

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
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	resp := toTagsResp(tags)
	c.JSON(http.StatusOK, resp)
}

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
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	c.JSON(http.StatusOK, toTagResp(tag))
}

type UpdateByIDTagReq struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}

// PATCH /tag/{id}
func (h *TagHandler) UpdateByID(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	var req UpdateByIDTagReq
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
	tag := &model.Tag{
		Name:   req.Name,
		ID:     id,
		UserID: userID,
	}
	updTag, err := h.s.UpdateByID(ctx, userID, tag)
	if err != nil {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	c.JSON(http.StatusOK, toTagResp(updTag))
}

func (h *TagHandler) DeleteByID(c *gin.Context) {
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
	c.JSON(http.StatusOK, "tag is delete")
}
