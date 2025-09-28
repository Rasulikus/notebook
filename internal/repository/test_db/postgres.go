package testdb

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"

	myfs "github.com/Rasulikus/notebook"
	"github.com/Rasulikus/notebook/internal/config"
	"github.com/Rasulikus/notebook/internal/model"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var (
	fs          embed.FS //nolint:unused // used
	testDB      *bun.DB
	testDSN     string
	truncateSQL = `
	TRUNCATE TABLE
		notes_tags,
		notes,
		tags,
		users
	RESTART IDENTITY CASCADE;
	`
)

func DB() *bun.DB {
	if testDB == nil {
		cfg := config.DbConfig{
			User: "admin",
			Pass: "password",
			Host: "localhost",
			Port: "5432",
			Name: "notebook_test",
		}

		var err error
		testDB, err = newClient(&cfg)
		if err != nil {
			log.Fatal(err)
		}
	}
	return testDB
}

func newClient(cfg *config.DbConfig) (*bun.DB, error) {
	testDSN = cfg.PostgresURL()
	// Open a PostgreSQL database
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(testDSN)))
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
	return db, nil
}

func CloseDB() {
	if testDB != nil {
		err := testDB.Close()
		if err != nil {
			log.Print(err)
			return
		}
		testDB = nil
	}
}

func RecreateTables() {
	if testDB == nil {
		DB()
	}
	d, err := iofs.New(myfs.Files, "migrations")
	if err != nil {
		log.Fatalf("iofs.New: %v", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, testDSN)
	if err != nil {
		log.Fatalf("migrate.NewWithSourceInstance: %v", err)
		return
	}
	defer m.Close() //nolint:errcheck // not need

	err = m.Force(1)
	if err != nil {
		log.Fatalf("migrate.Force: %v", err)
		return
	}
	if err := m.Down(); err != nil {
		log.Printf("migrate.Drop: %v", err)
	}

	if err := m.Up(); err != nil {
		log.Fatalf("migrate.Up: %v", err)
	}
}

func CleanDB(ctx context.Context) {
	_, err := testDB.ExecContext(ctx, truncateSQL)
	if err != nil {
		log.Fatalf("error clean db: %v", err)
	}
}
