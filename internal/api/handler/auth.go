// Package handler HTTP auth handlers (Gin).
package handler

import (
	"net/http"
	"time"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	refreshCookieName = "refresh_token"
	refreshCookiePath = "/auth"
)

// AuthHandler dependencies and settings for auth endpoints.
type AuthHandler struct {
	authService  service.AuthService
	refreshTTL   time.Duration
	secureCookie bool
}

// NewAuthHandler constructor.
func NewAuthHandler(authService service.AuthService, refreshTTL time.Duration, secureCookie bool) *AuthHandler {
	return &AuthHandler{authService: authService, refreshTTL: refreshTTL, secureCookie: secureCookie}
}

// RegisterReq registration request.
type RegisterReq struct {
	Email    string `json:"email"    binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Name     string `json:"name"     binding:"required,min=3,max=30"`
}

// Register (POST /auth/register) registers a user.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterReq
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
	if err := h.authService.Register(ctx, req.Email, req.Password, req.Name); err != nil {
		status, pub := model.ToHTTP(err)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	c.Status(http.StatusCreated)
}

// LoginReq login request.
type LoginReq struct {
	Email    string `json:"email"    binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// LoginResp login response.
type LoginResp struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
}

// Login (POST /auth/login) returns access token and user_id in JSON.
// Refresh token is stored in an HttpOnly cookie.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
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
	access, refresh, userID, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		status, pub := model.ToHTTP(err)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	setRefreshCookie(c, refresh, h.refreshTTL, h.secureCookie)
	c.JSON(http.StatusOK, LoginResp{AccessToken: access, UserID: userID})
}

// RefreshResp response with new access token.
type RefreshResp struct {
	AccessToken string `json:"access_token"`
}

// Refresh (POST /auth/refresh) reads refresh cookie, rotates refresh,
// updates the cookie, and returns a new access token in JSON.
func (h *AuthHandler) Refresh(c *gin.Context) {
	refresh, err := c.Cookie(refreshCookieName)
	if err != nil || refresh == "" {
		status, pub := model.ToHTTP(model.ErrUnauthorized)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	ctx := c.Request.Context()
	access, newRefresh, err := h.authService.Refresh(ctx, refresh)
	if err != nil {
		clearRefreshCookie(c, h.secureCookie)
		status, pub := model.ToHTTP(err)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	setRefreshCookie(c, newRefresh, h.refreshTTL, h.secureCookie)
	c.JSON(http.StatusOK, RefreshResp{AccessToken: access})
}

// Logout (POST /auth/logout) revokes the session and clears the refresh cookie.
func (h *AuthHandler) Logout(c *gin.Context) {
	refresh, err := c.Cookie(refreshCookieName)
	if err != nil || refresh == "" {
		clearRefreshCookie(c, h.secureCookie)
		status, pub := model.ToHTTP(model.ErrBadRequest)
		c.AbortWithStatusJSON(status, pub)
		return
	}
	ctx := c.Request.Context()
	if err := h.authService.Logout(ctx, refresh); err != nil {
		clearRefreshCookie(c, h.secureCookie)
		status, pub := model.ToHTTP(err)
		c.AbortWithStatusJSON(status, pub)
		return
	}

	clearRefreshCookie(c, h.secureCookie)
	c.Status(http.StatusNoContent)
}

// setRefreshCookie sets an HttpOnly refresh cookie (SameSite=Lax).
func setRefreshCookie(c *gin.Context, refresh string, ttl time.Duration, secureCookie bool) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookieName, refresh, int(ttl.Seconds()), refreshCookiePath, "", secureCookie, true)
}

// clearRefreshCookie removes the refresh cookie.
func clearRefreshCookie(c *gin.Context, secureCookie bool) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookieName, "", -1, refreshCookiePath, "", secureCookie, true)
}
