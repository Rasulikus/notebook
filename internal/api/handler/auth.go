package handler

import (
	"net/http"
	"time"

	"github.com/Rasulikus/notebook/internal/service"
	"github.com/gin-gonic/gin"
)

type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	Access string `json:"access"`
	UserID int64  `json:"user_id"`
}

type RefreshResp struct {
	Access string `json:"access"`
}

const (
	refreshCookieName = "refresh_token"
	refreshCookiePath = "/auth/refresh"
)

type AuthHandler struct {
	authService  service.AuthService
	jwtService   service.JWTService
	refreshTTL   time.Duration
	secureCookie bool
}

func NewAuthHanlder(authService service.AuthService, jwtService service.JWTService, refreshTTL time.Duration, secureCookie bool) *AuthHandler {
	return &AuthHandler{authService: authService, jwtService: jwtService, refreshTTL: refreshTTL, secureCookie: secureCookie}
}

// POST /auth/login
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	if err := h.authService.Register(ctx, req.Email, req.Password, req.Name); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	access, refresh, userID, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	setRefreshCookie(c, refresh, h.refreshTTL, h.secureCookie)
	c.JSON(http.StatusOK, LoginResp{
		Access: access,
		UserID: userID,
	})
}

// POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	refresh, err := c.Cookie(refreshCookieName)
	if err != nil || refresh == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	access, newRefresh, err := h.jwtService.RotateRefreshToken(ctx, refresh)
	if err != nil {
		//clearRefreshCookie(c, h.secureCookie)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	setRefreshCookie(c, newRefresh, h.refreshTTL, h.secureCookie)
	c.JSON(http.StatusOK, RefreshResp{
		Access: access,
	})

}

// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	refresh, err := c.Cookie(refreshCookieName)
	if err != nil || refresh == "" {
		clearRefreshCookie(c, h.secureCookie)
		c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	if err := h.jwtService.Logout(ctx, refresh); err != nil {
		clearRefreshCookie(c, h.secureCookie)
		c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	}

	clearRefreshCookie(c, h.secureCookie)
	c.Status(http.StatusNoContent)
}

func setRefreshCookie(c *gin.Context, refresh string, ttl time.Duration, secureCookie bool) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		refreshCookieName,
		refresh,
		int(ttl.Seconds()),
		refreshCookiePath,
		"",
		secureCookie,
		true,
	)
}

func clearRefreshCookie(c *gin.Context, secureCookie bool) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookieName, "", -1, refreshCookiePath, "", secureCookie, true)
}
