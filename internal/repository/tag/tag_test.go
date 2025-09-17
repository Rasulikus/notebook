package tag

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
	db      *bun.DB
	tagRepo *repo
	ctx     context.Context
}

func setupTestSuite(t *testing.T) *testSuite {
	t.Helper()
	var suite testSuite
	suite.db = testdb.DB()
	suite.tagRepo = NewRepository(suite.db)
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
	user := ensureUser(t, ts.db, ts.ctx)
	tests := []struct {
		name    string
		tag     *model.Tag
		wantErr bool
	}{
		{
			"create general tag",
			&model.Tag{
				Name: "base",
			},
			false,
		},
		{
			"create user tag",
			&model.Tag{
				Name:   "food",
				UserID: user.ID,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ts.tagRepo.Create(ts.ctx, tt.tag)
			if tt.wantErr {
				require.Error(t, err, tt.name)
				return
			}

			require.NoError(t, err)
			require.NotZero(t, tt.tag.ID)
		})
	}
}

func Test_Repo_List(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	err := ts.tagRepo.Create(ts.ctx, &model.Tag{
		Name: "base1",
	})
	require.NoError(t, err)
	err = ts.tagRepo.Create(ts.ctx, &model.Tag{
		Name:   "base2",
		UserID: user.ID,
	})
	require.NoError(t, err)

	tags, err := ts.tagRepo.List(ts.ctx, user.ID)
	require.NoError(t, err)

	require.Equal(t, 2, len(tags))
}
