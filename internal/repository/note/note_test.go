package note

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository/tag"
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
	tagRepo  *tag.Repo
	ctx      context.Context
}

func setupTestSuite(t *testing.T) *testSuite {
	t.Helper()
	var suite testSuite
	suite.db = testdb.DB()
	suite.noteRepo = NewRepository(suite.db)
	suite.tagRepo = tag.NewRepository(suite.db)
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

	cases := []struct {
		name      string
		buildNote func(userID int64) *model.Note
		buildTags func(userID int64) []*model.Tag
		wantErr   bool
		check     func(t *testing.T, n *model.Note)
	}{
		{
			name: "ok: create without tags",
			buildNote: func(userID int64) *model.Note {
				return &model.Note{
					Title:  "test note",
					Text:   "my no error test note",
					UserID: userID,
				}
			},
			buildTags: func(userID int64) []*model.Tag {
				return nil
			},
			wantErr: false,
			check: func(t *testing.T, n *model.Note) {
				require.NotZero(t, n.ID)
				require.Equal(t, "test note", n.Title)
				require.Equal(t, "my no error test note", n.Text)
			},
		},
		{
			name: "ok: create with tags",
			buildNote: func(userID int64) *model.Note {
				return &model.Note{
					Title:  "Note",
					Text:   "note Text",
					UserID: userID,
				}
			},
			buildTags: func(userID int64) []*model.Tag {
				t1 := &model.Tag{Name: "Sport", UserID: userID}
				t2 := &model.Tag{Name: "Chicken", UserID: userID}
				require.NoError(t, ts.tagRepo.Create(ts.ctx, t1))
				require.NoError(t, ts.tagRepo.Create(ts.ctx, t2))
				return []*model.Tag{t1, t2}
			},
			wantErr: false,
			check: func(t *testing.T, n *model.Note) {
				require.NotZero(t, n.ID)
			},
		},
		{
			name: "error: invalid user",
			buildNote: func(_ int64) *model.Note {
				return &model.Note{
					Title:  "test note",
					Text:   "my error test note",
					UserID: 999999999, // несуществующий пользователь
				}
			},
			buildTags: func(userID int64) []*model.Tag { return nil },
			wantErr:   true,
			check:     func(t *testing.T, n *model.Note) {},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testdb.CleanDB(ts.ctx)
			user := ensureUser(t, ts.db, ts.ctx)

			note := tc.buildNote(user.ID)
			tags := tc.buildTags(user.ID)

			newNote, err := ts.noteRepo.Create(ts.ctx, note, tags)
			if tc.wantErr {
				require.Error(t, err, tc.name)
				return
			}
			require.NoError(t, err)
			tc.check(t, newNote)
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
	require.ErrorIs(t, err, model.ErrNotFound, "there is no note with this id")
	require.Nil(t, gotErr)

	gotErr2, err := ts.noteRepo.GetByID(ts.ctx, 9999999, n.ID)
	require.ErrorIs(t, err, model.ErrNotFound, "there is no note with this user_id")
	require.Nil(t, gotErr2)
}

func Test_Repo_Update(t *testing.T) {
	now := time.Now().UTC()
	log.Println(now)
	updTitle := "updated title"
	updText := "updated text"
	tags := []*model.Tag{{Name: "tag1"}, {Name: "tag2"}}
	updTags := []*model.Tag{{Name: "newTag1"}}
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
	_, err := ts.tagRepo.CreateTags(ts.ctx, tags)
	require.NoError(t, err)

	_, err = ts.tagRepo.CreateTags(ts.ctx, updTags)
	require.NoError(t, err)

	newNote, err := ts.noteRepo.Create(ts.ctx, n, tags)
	require.NoError(t, err, "create note error")
	require.WithinDuration(t, hourEarlie, newNote.UpdatedAt, time.Second)
	require.Len(t, newNote.Tags, len(tags))

	updNote, err := ts.noteRepo.UpdateByID(ts.ctx, n.UserID, n.ID, &updTitle, &updText, &[]int64{updTags[0].ID})
	require.NoError(t, err, "update note error")
	require.Len(t, updNote.Tags, len(updTags))
	require.Equal(t, updTags[0].ID, updNote.Tags[0].ID)
	require.WithinDuration(t, now, updNote.UpdatedAt, time.Second)
}

func Test_Repo_DeleteByID(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	n := insertNote(t, ts.db, ts.ctx, user.ID, "n1", "note 1")

	err := ts.noteRepo.DeleteByID(ts.ctx, n.UserID, n.ID)
	require.NoError(t, err)

	got, err := ts.noteRepo.GetByID(ts.ctx, n.UserID, n.ID)
	require.ErrorIs(t, err, model.ErrNotFound, "there is no note with this id")
	require.Nil(t, got)
}

func Test_Repo_DeleteTags(t *testing.T) {
	ts := setupTestSuite(t)
	testdb.CleanDB(ts.ctx)
	user := ensureUser(t, ts.db, ts.ctx)
	tags := []*model.Tag{{Name: "sport"}, {Name: "football"}}
	_, err := ts.tagRepo.CreateTags(ts.ctx, tags)
	require.NoError(t, err)
	note := &model.Note{Title: "title", Text: "text", UserID: user.ID}
	n, err := ts.noteRepo.Create(ts.ctx, note, tags)
	require.NoError(t, err)
	require.Len(t, n.Tags, len(tags))
	err = ts.noteRepo.DeleteTags(ts.ctx, user.ID, n.ID)
	require.NoError(t, err)
	var cnt int
	err = ts.db.NewSelect().Table("notes_tags").ColumnExpr("count(*)").Where("note_id = ?", n.ID).Scan(ts.ctx, &cnt)
	require.NoError(t, err)
	require.Equal(t, 0, cnt)
}
