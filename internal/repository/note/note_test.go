package note

import (
	"context"
	"database/sql"
	"log"
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
	db       *bun.DB
	noteRepo *repo
	ctx      context.Context
}

func setupTestSuite(t *testing.T) *testSuite {
	t.Helper()
	var suite testSuite
	suite.db = testdb.DB()
	suite.noteRepo = NewRepository(suite.db)
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

func insertNote(t *testing.T, db *bun.DB, ctx context.Context, userID int64, title, text string) *model.Note {
	t.Helper()
	n := &model.Note{
		Title:  title,
		Text:   text,
		UserID: userID,
	}
	err := db.NewInsert().Model(n).Scan(ctx, n)
	require.NoError(t, err)
	require.NotZero(t, n.ID)
	return n
}

func Test_Repo_Create(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	tests := []struct {
		name    string
		note    *model.Note
		wantErr bool
	}{
		{
			"create no err note",
			&model.Note{
				Title:  "test note",
				Text:   "my no error test note",
				UserID: user.ID,
			},
			false,
		},
		{
			"invalid user",
			&model.Note{
				Title:  "test note",
				Text:   "my error test note",
				UserID: 999999999,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ts.noteRepo.Create(ts.ctx, tt.note)
			if tt.wantErr {
				require.Error(t, err, tt.name)
				return
			}

			require.NoError(t, err)
			require.NotZero(t, tt.note.ID)
		})
	}
}

func Test_Repo_List(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	insertNote(t, ts.db, ts.ctx, user.ID, "n1", "note 1")
	insertNote(t, ts.db, ts.ctx, user.ID, "n2", "note 2")
	insertNote(t, ts.db, ts.ctx, user.ID, "n3", "note 3")

	list, err := ts.noteRepo.List(ts.ctx, user.ID, 20, 0, "")
	require.NoError(t, err)
	require.Equal(t, len(list), 3)

	listWithLimit, err := ts.noteRepo.List(ts.ctx, user.ID, 1, 0, "")
	require.NoError(t, err)
	require.Equal(t, len(listWithLimit), 1)

	listWithOffset, err := ts.noteRepo.List(ts.ctx, user.ID, 20, 1, "")
	require.NoError(t, err)
	require.Equal(t, len(listWithOffset), 2)
}

func Test_Repo_GetByID(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	n := insertNote(t, ts.db, ts.ctx, user.ID, "n1", "note 1")

	got, err := ts.noteRepo.GetByID(ts.ctx, n.UserID, n.ID)
	require.NoError(t, err)
	require.Equal(t, got.ID, n.ID)

	gotErr, err := ts.noteRepo.GetByID(ts.ctx, n.UserID, 99999999)
	require.ErrorIs(t, err, sql.ErrNoRows, "there is no note with this id")
	require.Nil(t, gotErr)

	gotErr2, err := ts.noteRepo.GetByID(ts.ctx, 9999999, n.ID)
	require.ErrorIs(t, err, sql.ErrNoRows, "there is no note with this user_id")
	require.Nil(t, gotErr2)
}

func Test_Repo_Update(t *testing.T) {
	now := time.Now().UTC()
	log.Println(now)
	updTitle := "updated title"
	updText := "updated text"
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	hourEarlie := now.Add(-1 * time.Hour)
	n := &model.Note{
		Title:     "note",
		Text:      "my note",
		UserID:    user.ID,
		CreatedAt: hourEarlie,
		UpdatedAt: hourEarlie,
	}
	err := ts.noteRepo.Create(ts.ctx, n)
	require.WithinDuration(t, hourEarlie, n.UpdatedAt, time.Second)
	require.NoError(t, err, "create note error")
	n.Title = updTitle
	n.Text = updText
	err = ts.noteRepo.UpdateByID(ts.ctx, n.UserID, n)
	require.NoError(t, err, "update note error")
	require.WithinDuration(t, now, n.UpdatedAt, time.Second)
	require.Equal(t, updTitle, n.Title)
	require.Equal(t, updText, n.Text)
}

func Test_Repo_DeleteByID(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	n := insertNote(t, ts.db, ts.ctx, user.ID, "n1", "note 1")

	err := ts.noteRepo.DeleteByID(ts.ctx, n.UserID, n.ID)
	require.NoError(t, err)

	got, err := ts.noteRepo.GetByID(ts.ctx, n.UserID, n.ID)
	require.ErrorIs(t, err, sql.ErrNoRows, "there is no note with this id")
	require.Nil(t, got)
}
