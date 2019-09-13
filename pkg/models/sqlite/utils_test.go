package sqlite_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", "file:memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	// DB migrations
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// Look for migration files relative to this directory
	filePath, err := os.Getwd()
	migrationsPath := fmt.Sprintf("%s/../../../sql/migrations", filePath)
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://"+migrationsPath),
		"sqlite3", driver)
	if err != nil {
		t.Fatal(err)
	}

	err = m.Up()

	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		err := m.Down()
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}
