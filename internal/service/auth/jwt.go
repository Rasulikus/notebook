package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bun/driver/pgdriver"
)

const maxTries = 3

type TokenManager struct {
	secret      []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
	sessionRepo repository.SessionRepository
}

type Claims struct {
	jwt.RegisteredClaims
}

func NewTokenManager(secret []byte, accessTTL, refreshTTL time.Duration, sessionRepo repository.SessionRepository) *TokenManager {
	return &TokenManager{
		secret:      secret,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
		sessionRepo: sessionRepo,
	}
}

func (m *TokenManager) NewAccessToken(userID int64) (string, error) {
	now := time.Now().UTC()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func generateRefreshToken() (string, []byte, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(token))

	return token, sum[:], nil
}

func (m *TokenManager) CreateRefreshToken(ctx context.Context, userID int64) (string, error) {
	now := time.Now().UTC()

	for i := 0; i < maxTries; i++ {
		token, hash, err := generateRefreshToken()
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
			if isUniqueViolation(err) {
				continue
			}
			return "", err
		}
		return token, nil
	}
	return "", model.ErrInvalidToken
}

func (m *TokenManager) UpdateRefreshToken(ctx context.Context, oldRefresh string) (string, string, error) {
	now := time.Now().UTC()
	oldhash := sha256.Sum256([]byte(oldRefresh))

	for i := 0; i < maxTries; i++ {
		newRefresh, newHash, err := generateRefreshToken()
		if err != nil {
			return "", "", err
		}
		newSession, err := m.sessionRepo.RotateRefreshTokenTx(ctx, oldhash[:], newHash, now.Add(m.refreshTTL))
		if err != nil {
			if isUniqueViolation(err) {
				continue
			}
			return "", "", err
		}
		access, err := m.NewAccessToken(newSession.UserID)
		if err != nil {
			return "", "", err
		}
		return access, newRefresh, nil
	}
	return "", "", model.ErrInvalidToken
}

func (m *TokenManager) ParseAccessToken(token string) (int64, error) {
	var claims Claims

	t, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, model.ErrInvalidToken
			}
			return m.secret, nil
		})
	if err != nil || !t.Valid {
		return 0, model.ErrInvalidToken
	}
	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, model.ErrInvalidToken
	}
	return userID, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgdriver.Error
	if errors.As(err, &pgErr) && pgErr.Field('C') == "23505" {
		return true
	}
	return false
}
