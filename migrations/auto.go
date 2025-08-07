package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/Rasulikus/notebook/internal/config"
	"github.com/Rasulikus/notebook/internal/repository"
	"github.com/pressly/goose/v3"
)

func main() {

	var cmd string
	flag.StringVar(&cmd, "cmd", "up", "goose command: up, down, status")
	flag.Parse()

	cfg := config.LoadConfig()

	bunDB, err := repository.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := bunDB.DB.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()

	sqlDB := bunDB.DB.DB

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	migDir := filepath.Join(wd, "migrations")

	switch cmd {
	case "up":
		if err := goose.Up(sqlDB, migDir); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := goose.Down(sqlDB, migDir); err != nil {
			log.Fatal(err)
		}
	case "status":
		if err := goose.Status(sqlDB, migDir); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown cmd %s", cmd)
	}
}
