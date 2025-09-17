package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret      []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
	sessionRepo repository.SessionRepository
}

func NewTokenManager(secret []byte, accessTTL, refreshTTL time.Duration, sessionRepo repository.SessionRepository) *JWTService {
	return &JWTService{
		secret:      secret,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
		sessionRepo: sessionRepo,
	}
}

func (m *JWTService) CreateAccessToken(userID int64) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"uid": userID,
		"exp": now.Add(m.accessTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (m *JWTService) CreateRefreshToken(ctx context.Context, userID int64) (string, error) {
	now := time.Now().UTC()

	token, hash, err := generateRefreshTokenWithHash()
	if err != nil {
		return "", err
	}

	sess := &model.Session{
		UserID:           userID,
		RefreshTokenHash: hash,
		ExpiresAt:        now.Add(m.refreshTTL),
	}

	err = m.sessionRepo.Create(ctx, sess)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *JWTService) RotateRefreshToken(ctx context.Context, oldRefresh string) (string, string, error) {
	now := time.Now().UTC()
	oldhash := generateRefreshTokenHash(oldRefresh)

	newRefresh, newHash, err := generateRefreshTokenWithHash()
	if err != nil {
		return "", "", err
	}
	newSession, err := m.sessionRepo.RotateRefreshTokenTx(ctx, oldhash[:], newHash, now.Add(m.refreshTTL))
	if err != nil {
		return "", "", err
	}
	access, err := m.CreateAccessToken(newSession.UserID)
	if err != nil {
		return "", "", err
	}
	return access, newRefresh, nil
}

func (m *JWTService) ParseAccessToken(token string) (int64, error) {
	var claims jwt.MapClaims

	t, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, model.ErrBadRequest
			}
			return m.secret, nil
		})
	if err != nil || !t.Valid {
		return 0, model.ErrBadRequest
	}
	userID, ok := claims["uid"].(float64)
	if !ok {
		return 0, model.ErrBadRequest
	}
	return int64(userID), nil
}

func generateRefreshTokenWithHash() (string, []byte, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(token))

	return token, sum[:], nil
}

func generateRefreshTokenHash(refreshToken string) []byte {
	refreshTokenHash := sha256.Sum256([]byte(refreshToken))
	return refreshTokenHash[:]
}
