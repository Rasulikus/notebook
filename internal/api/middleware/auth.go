// Package middleware - auth-related Gin middleware and helpers.
package middleware

import (
	"strings"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/service"
	"github.com/gin-gonic/gin"
)

// ctxUserIDKey - context key for storing authenticated user ID.
const ctxUserIDKey = "userID"

// AuthMiddleware validates "Authorization: Bearer <access>" header,
// parses the access token, and stores user ID in context.
// Aborts the request with a public error on failure.
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(authz, prefix) {
			status, pub := model.ToHTTP(model.ErrUnauthorized)
			c.AbortWithStatusJSON(status, pub)
			return
		}
		tokenStr := strings.TrimPrefix(authz, prefix)

		uid, err := authService.ParseAccessToken(tokenStr)
		if err != nil {
			status, pub := model.ToHTTP(err)
			c.AbortWithStatusJSON(status, pub)
			return
		}

		c.Set(ctxUserIDKey, uid)
		c.Next()
	}
}

// CurrentUserID extracts user ID from context.
// Aborts with an error if the ID is missing or of a wrong type.
func CurrentUserID(c *gin.Context) int64 {
	v, ok := c.Get(ctxUserIDKey)
	if !ok {
		status, pub := model.ToHTTP(model.ErrUnauthorized)
		c.AbortWithStatusJSON(status, pub)
		return 0
	}
	id, ok := v.(int64)
	if !ok {
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return 0
	}
	return id
}
