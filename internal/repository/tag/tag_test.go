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
	tagRepo *Repo
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

func Test_Repo_CreateTags(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	tags := []*model.Tag{{Name: "sport", UserID: user.ID}, {Name: "sport", UserID: user.ID}, {Name: "football"}}
	got, err := ts.tagRepo.CreateTags(ts.ctx, tags)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.NotZero(t, got[0].ID)
	require.NotZero(t, got[1].ID)
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

	tags, err := ts.tagRepo.List(ts.ctx, user.ID, 10, 0, "")
	require.NoError(t, err)

	require.Equal(t, 2, len(tags))
}

func Test_Repo_GetByID(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	tag := &model.Tag{Name: "tag1", UserID: user.ID}
	err := ts.tagRepo.Create(ts.ctx, tag)
	require.NoError(t, err)

	got, err := ts.tagRepo.GetByID(ts.ctx, tag.UserID, tag.ID)
	require.NoError(t, err)
	require.Equal(t, got.ID, tag.ID)

	gotErr, err := ts.tagRepo.GetByID(ts.ctx, tag.UserID, 99999999)
	require.ErrorIs(t, err, model.ErrNotFound, "there is no tag with this id")
	require.Nil(t, gotErr)

	gotErr2, err := ts.tagRepo.GetByID(ts.ctx, 9999999, tag.ID)
	require.ErrorIs(t, err, model.ErrNotFound, "there is no tag with this user_id")
	require.Nil(t, gotErr2)
}

func Test_Repo_GetByIDs(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)

	u1 := ensureUser(t, ts.db, ts.ctx)

	t1 := &model.Tag{Name: "sport", UserID: u1.ID}
	t2 := &model.Tag{Name: "travel", UserID: u1.ID}
	t3 := &model.Tag{Name: "sport", UserID: u1.ID}

	_, err := ts.tagRepo.CreateTags(ts.ctx, []*model.Tag{t1, t2, t3})
	require.NoError(t, err)

	got, err := ts.tagRepo.GetByIDs(ts.ctx, u1.ID, []int64{t1.ID, t2.ID})
	require.NoError(t, err)
	require.Len(t, got, 2)
	ids := map[int64]bool{}
	for _, tg := range got {
		ids[tg.ID] = true
		require.Equal(t, u1.ID, tg.UserID)
	}
	require.True(t, ids[t1.ID])
	require.True(t, ids[t2.ID])

	got, err = ts.tagRepo.GetByIDs(ts.ctx, u1.ID, []int64{t1.ID, t2.ID, t3.ID})
	require.Error(t, err)
	require.Nil(t, got)

	got, err = ts.tagRepo.GetByIDs(ts.ctx, u1.ID, nil)
	require.NoError(t, err)
	require.Nil(t, got)

	got, err = ts.tagRepo.GetByIDs(ts.ctx, u1.ID, []int64{999999, 888888})
	require.Error(t, err)
	require.Nil(t, got)
}

func Test_Repo_UpdateByID(t *testing.T) {
	newTagName := "updated Name"
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	tag := &model.Tag{Name: "tag", UserID: user.ID}
	err := ts.tagRepo.Create(ts.ctx, tag)
	require.NoError(t, err)
	tag.Name = newTagName
	updTag, err := ts.tagRepo.UpdateByID(ts.ctx, user.ID, tag)
	require.NoError(t, err, "update tag error")
	require.Equal(t, newTagName, updTag.Name)

	errTag, err := ts.tagRepo.UpdateByID(ts.ctx, 9999999, tag)
	require.Error(t, err, "no tag with that user id")
	require.Empty(t, errTag)
}
func Test_Repo_DeleteByID(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	tag := &model.Tag{Name: "tag1", UserID: user.ID}
	err := ts.tagRepo.Create(ts.ctx, tag)
	require.NoError(t, err)
	err = ts.tagRepo.DeleteByID(ts.ctx, user.ID, tag.ID)
	require.NoError(t, err)

	got, err := ts.tagRepo.GetByID(ts.ctx, user.ID, tag.ID)
	require.ErrorIs(t, err, model.ErrNotFound, "there is no tag with this id")
	require.Nil(t, got)
}
