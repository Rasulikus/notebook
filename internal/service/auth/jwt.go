package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

// tokenManager issues/parses access tokens and manages refresh sessions.
type tokenManager struct {
	TokenConfig
}

// NewTokenManager constructor for JWTService.
func newTokenManager(cfg TokenConfig) *tokenManager {
	return &tokenManager{cfg}
}

// CreateAccessToken builds HS256 JWT with uid and exp.
func (m *tokenManager) CreateAccessToken(userID int64) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"uid": userID,
		"exp": now.Add(m.AccessTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.Secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

// generateRefreshTokenWithHash creates a random refresh token (base64url)
// and returns its HMAC-SHA256 (with service secret). Only the HMAC is stored in DB.
func (m *tokenManager) generateRefreshTokenWithHash() (string, []byte, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil, err
	}
	token := base64.RawURLEncoding.EncodeToString(b)

	mac := hmac.New(sha256.New, m.Secret)
	mac.Write([]byte(token))
	sum := mac.Sum(nil)

	return token, sum, nil
}

// generateRefreshTokenHash computes HMAC for a provided refresh token string.
func (m *tokenManager) generateRefreshTokenHash(refreshToken string) []byte {
	mac := hmac.New(sha256.New, m.Secret)
	mac.Write([]byte(refreshToken))
	return mac.Sum(nil)
}

// CreateRefreshToken creates a refresh token and persists its HMAC with TTL.
func (m *tokenManager) CreateRefreshToken(ctx context.Context, userID int64) (string, error) {
	now := time.Now().UTC()

	token, hash, err := m.generateRefreshTokenWithHash()
	if err != nil {
		return "", err
	}

	sess := &model.Session{
		UserID:           userID,
		RefreshTokenHash: hash,
		ExpiresAt:        now.Add(m.RefreshTTL),
	}

	err = m.SessionRepo.Create(ctx, sess)
	if err != nil {
		return "", err
	}
	return token, nil
}

// RotateRefreshToken atomically rotates a valid, unrevoked, unexpired refresh.
func (m *tokenManager) RotateRefreshToken(ctx context.Context, oldRefresh string) (string, string, error) {
	now := time.Now().UTC()
	oldhash := m.generateRefreshTokenHash(oldRefresh)

	newRefresh, newHash, err := m.generateRefreshTokenWithHash()
	if err != nil {
		return "", "", err
	}
	newSession, err := m.SessionRepo.RotateRefreshToken(ctx, oldhash[:], newHash, now.Add(m.RefreshTTL))
	if err != nil {
		return "", "", err
	}
	access, err := m.CreateAccessToken(newSession.UserID)
	if err != nil {
		return "", "", err
	}
	return access, newRefresh, nil
}

// ParseAccessToken verifies HS256, validates signature/exp,
// and extracts uid claim as int64.
func (m *tokenManager) ParseAccessToken(token string) (int64, error) {
	var claims jwt.MapClaims

	t, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, model.ErrBadRequest
			}
			return m.Secret, nil
		})
	if err != nil || !t.Valid {
		return 0, model.ErrUnauthorized
	}
	userID, ok := claims["uid"].(float64)
	if !ok {
		return 0, model.ErrUnauthorized
	}
	return int64(userID), nil
}
