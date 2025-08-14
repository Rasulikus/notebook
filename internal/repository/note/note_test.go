package note

import (
	"context"
	"os"
	"testing"

	"github.com/Rasulikus/notebook/internal/model"
	testdb "github.com/Rasulikus/notebook/internal/repository/test_db"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

var note *model.Note

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
	u := &model.User{Email: "test@mail.ru", PasswordHash: "x", UserName: "test"}
	_, err := db.NewInsert().Model(u).Exec(ctx)
	require.NoError(t, err)
	return u
}

func Test_repo_Create(t *testing.T) {
	ts := setupTestSuite(t)

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
				UserID: 5,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ts.noteRepo.Create(ts.ctx, tt.note)
			if tt.wantErr {
				require.Error(t, err, "there is no user with this id")
				return
			}

			require.NoError(t, err)
			require.NotZero(t, tt.note.ID)
		})
	}
}
