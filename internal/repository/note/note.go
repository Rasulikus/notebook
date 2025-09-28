package note

import (
	"context"
	"fmt"

	"github.com/Rasulikus/notebook/internal/model"
	"github.com/Rasulikus/notebook/internal/repository"
	"github.com/uptrace/bun"
)

type repo struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *repo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, note *model.Note, tags []*model.Tag) (*model.Note, error) {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(note).Returning("*").Exec(ctx)
		if err != nil {
			return err
		}
		if len(tags) == 0 {
			return nil
		}
		for _, tag := range tags {
			if tag == nil || tag.ID == 0 {
				return fmt.Errorf("invalid tag: nil or id=0")
			}
			_, err = tx.NewInsert().
				Model(&model.NoteTag{NoteID: note.ID, TagID: tag.ID}).
				Column("note_id", "tag_id").
				On("CONFLICT (note_id, tag_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, note.UserID, note.ID)
}

func (r *repo) List(ctx context.Context, userID int64, limit, offset int, order string) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.NewSelect().
		Model(&notes).
		Where("user_id = ?", userID).
		Relation("Tags").
		Order(order).
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *repo) GetByID(ctx context.Context, userID, id int64) (*model.Note, error) {
	note := new(model.Note)
	err := r.db.NewSelect().Model(note).Relation("Tags").Where("id = ? AND user_id = ?", id, userID).Scan(ctx)
	if err != nil {
		return nil, repository.IsNoRowsError(err)
	}
	return note, nil
}

func (r *repo) UpdateByID(ctx context.Context, userID, id int64, title, text *string, tagsIDs *[]int64) (*model.Note, error) {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		q := tx.NewUpdate().
			Model((*model.Note)(nil)).
			Set("updated_at = now()").
			Where("id = ? AND user_id = ?", id, userID)
		if title != nil {
			q.Set("title = ?", *title)
		}
		if text != nil {
			q.Set("text = ?", *text)
		}
		res, err := q.Exec(ctx)
		if err != nil {
			return err
		}

		aff, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if aff == 0 {
			return model.ErrNotFound
		}
		if tagsIDs != nil {
			_, err = tx.NewDelete().
				Table("notes_tags").
				Where("note_id = ?", id).
				Exec(ctx)
			if err != nil {
				return err
			}
			for _, tagID := range *tagsIDs {
				if tagID == 0 {
					return fmt.Errorf("invalid tag")
				}

				_, err = tx.NewInsert().
					Model(&model.NoteTag{NoteID: id, TagID: tagID}).
					Column("note_id", "tag_id").
					On("CONFLICT (note_id, tag_id) DO NOTHING").
					Exec(ctx)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	note, err := r.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (r *repo) DeleteByID(ctx context.Context, userID, id int64) error {
	res, err := r.db.NewDelete().Model((*model.Note)(nil)).Where("id = ? AND user_id = ?", id, userID).Exec(ctx)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *repo) DeleteTags(ctx context.Context, userID, id int64) error {
	res, err := r.db.NewDelete().
		TableExpr("notes_tags AS nt").
		Where("nt.note_id = ?", id).
		Where("EXISTS (SELECT 1 FROM notes n WHERE n.id = nt.note_id AND n.user_id = ?)", userID).
		Exec(ctx)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return model.ErrNotFound
	}
	return nil
}
