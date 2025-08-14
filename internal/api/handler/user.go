package handler

import "github.com/gin-gonic/gin"

func (h *NoteHandler) RegisterUsers(r *gin.RouterGroup) {
	r.POST("/", h.create)
	r.GET("/", h.list)
}
