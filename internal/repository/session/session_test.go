package session

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Rasulikus/notebook/internal/model"
	testdb "github.com/Rasulikus/notebook/internal/repository/test_db"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestMain(m *testing.M) {
	testdb.RecreateTables()
	code := m.Run()
	testdb.CloseDB()
	os.Exit(code)
}

type testSuite struct {
	db          *bun.DB
	sessionRepo *repo
	ctx         context.Context
}

func setupTestSuite(t *testing.T) *testSuite {
	t.Helper()
	var suite testSuite
	suite.db = testdb.DB()
	suite.sessionRepo = NewRepository(suite.db)
	suite.ctx = context.Background()
	return &suite
}

func ensureUser(t *testing.T, db *bun.DB, ctx context.Context) *model.User {
	t.Helper()
	u := &model.User{Email: "test@mail.ru", PasswordHash: "x", Name: "test"}
	err := db.NewInsert().Model(u).Scan(ctx, u)
	require.NoError(t, err)
	require.NotZero(t, u.ID)
	return u
}

func Test_Repo_Create(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	now := time.Now().UTC()
	user := ensureUser(t, ts.db, ts.ctx)
	tests := []struct {
		name    string
		session *model.Session
		wantErr bool
	}{
		{
			"create no err session",
			&model.Session{
				RefreshTokenHash: []byte("123456"),
				ExpiresAt:        now.Add(time.Hour),
				UserID:           user.ID,
			},
			false,
		},
		{
			"create no err session",
			&model.Session{
				RefreshTokenHash: []byte("123123"),
				UserID:           user.ID,
				ExpiresAt:        now.Add(time.Hour),
				RevokedAt:        time.Time{},
			},
			false,
		},
		{
			"invalid userID",
			&model.Session{
				RefreshTokenHash: []byte("123"),
				ExpiresAt:        now.Add(time.Hour),
				UserID:           999999999,
			},
			true,
		},
		{
			"not unique refresh token hash",
			&model.Session{
				RefreshTokenHash: []byte("123123"),
				ExpiresAt:        now.Add(time.Hour),
				UserID:           user.ID,
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ts.sessionRepo.Create(ts.ctx, tt.session)
			if tt.wantErr {
				require.Error(t, err, tt.name)
				return
			}
			require.NoError(t, err)
			require.NotZero(t, tt.session.ID)
		})
	}
}

func Test_Repo_RotateRefreshToken(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	now := time.Now().UTC()
	user := ensureUser(t, ts.db, ts.ctx)
	newRefreshTokenHash := []byte("654321")
	session := &model.Session{
		RefreshTokenHash: []byte("123456"),
		ExpiresAt:        now.Add(time.Hour),
		UserID:           user.ID,
	}

	err := ts.sessionRepo.Create(ts.ctx, session)
	require.NoError(t, err)
	newSession, err := ts.sessionRepo.RotateRefreshToken(ts.ctx, session.RefreshTokenHash, newRefreshTokenHash, now.Add(4*time.Hour))
	require.NoError(t, err)
	require.Equal(t, newSession.ID, session.ID)
	require.Equal(t, newSession.RefreshTokenHash, newRefreshTokenHash)
}

func Test_Repo_SetRevokedAtNow(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	now := time.Now().UTC()
	user := ensureUser(t, ts.db, ts.ctx)
	session := &model.Session{
		RefreshTokenHash: []byte("123456"),
		ExpiresAt:        now.Add(time.Hour),
		UserID:           user.ID,
	}
	err := ts.sessionRepo.Create(ts.ctx, session)
	require.NoError(t, err)

	err = ts.sessionRepo.SetRevokedAtNow(ts.ctx, session.RefreshTokenHash)
	require.NoError(t, err)

	err = ts.sessionRepo.SetRevokedAtNow(ts.ctx, session.RefreshTokenHash)
	require.Error(t, err, "error because revorkedAt already have")
}
