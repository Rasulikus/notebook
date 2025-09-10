package middleware

import (
	"net/http"
	"strings"

	"github.com/Rasulikus/notebook/internal/service"
	"github.com/gin-gonic/gin"
)

const ctxUserIDKey = "userID"

func AuthMiddleware(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(authz, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenStr := strings.TrimPrefix(authz, prefix)
		uid, err := jwtService.ParseAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}
		c.Set(ctxUserIDKey, uid)
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(ctxUserIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}
