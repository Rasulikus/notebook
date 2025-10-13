package user

import (
	"context"
	"os"
	"testing"

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
	db       *bun.DB
	userRepo *repo
	ctx      context.Context
}

func setupTestSuite(t *testing.T) *testSuite {
	t.Helper()
	var suite testSuite
	suite.db = testdb.DB()
	suite.userRepo = NewRepository(suite.db)
	suite.ctx = context.Background()
	return &suite
}

func Test_Repo_Create(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := &model.User{
		Email:        "abvg@g.ru",
		PasswordHash: "123",
		Name:         "ab",
	}
	err := ts.userRepo.Create(ts.ctx, user)
	require.NoError(t, err)
	require.NotZero(t, user.ID)
	err = ts.userRepo.Create(ts.ctx, user)
	require.ErrorIs(t, err, model.ErrConflict, "the user with this email already exists")
}

func Test_Repo_FindByEmail(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := &model.User{
		Email:        "abvg@g.ru",
		PasswordHash: "123",
		Name:         "ab",
	}
	err := ts.db.NewInsert().Model(user).Scan(ts.ctx, user)
	require.NoError(t, err)
	got, err := ts.userRepo.GetByEmail(ts.ctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, got, user)
	gotErr, err := ts.userRepo.GetByEmail(ts.ctx, "afdfffd@tetet.cl")
	require.Error(t, err)
	require.Equal(t, gotErr, (*model.User)(nil))
}
