package middleware

import (
	"strings"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/service"
	"github.com/gin-gonic/gin"
)

const ctxUserIDKey = "userID"

func AuthMiddleware(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(authz, prefix) {
			status, pub := model.ToHTTP(model.ErrUnauthorized)
			c.AbortWithStatusJSON(status, pub)
			return
		}
		tokenStr := strings.TrimPrefix(authz, prefix)
		uid, err := jwtService.ParseAccessToken(tokenStr)
		if err != nil {
			status, pub := model.ToHTTP(model.ErrUnauthorized)
			c.AbortWithStatusJSON(status, pub)
			return
		}
		c.Set(ctxUserIDKey, uid)
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) int64 {
	v, ok := c.Get(ctxUserIDKey)
	if !ok {
		status, pub := model.ToHTTP(model.ErrUnauthorized)
		c.AbortWithStatusJSON(status, pub)
		return 0
	}
	id, ok := v.(int64)
	if !ok {
		status, pub := model.ToHTTP(model.ErrConflict)
		c.AbortWithStatusJSON(status, pub)
		return 0
	}
	return id
}
