package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Rasulikus/notebook/internal/config"
	"github.com/Rasulikus/notebook/internal/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type DB struct {
	DB *bun.DB
}

func NewClient(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Db.User, cfg.Db.Pass, cfg.Db.Host, cfg.Db.Port, cfg.Db.Name)
	// Open a PostgreSQL database
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	err := sqldb.Ping()
	if err != nil {
		return nil, fmt.Errorf("cant connect to database: %w", err)
	}
	// Open a PostgreSQL database
	db := bun.NewDB(sqldb, pgdialect.New())
	// Print all queries to stdout
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	// add m2m
	db.RegisterModel((*model.NoteTag)(nil))
	return &DB{
		DB: db,
	}, nil
}

func IsNoRowsError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	return err
}
