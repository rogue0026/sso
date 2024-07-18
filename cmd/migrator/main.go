package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	storagePath     string
	migrationsPath  string
	migrationsTable string
)

func main() {
	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to dir with migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "", "migrations table name")
	flag.Parse()
	if storagePath == "" || migrationsPath == "" || migrationsTable == "" {
		panic("storage-path or migrations-path or migrations table name is not defined")
	}
	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("sqlite://%s?x-migrations-table=%s", storagePath, migrationsTable))
	if err != nil {
		panic(err.Error())
	}

	err = migrator.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err.Error())
	}
}
