package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// DB connection
	db, err := sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		errorLog.Fatal(err)
	} else {
		infoLog.Println("Opened DB...")
	}

	if err = db.Ping(); err != nil {
		errorLog.Fatal("No DB present")
	}

	defer db.Close()

	// DB migrations
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		errorLog.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./sql/migrations",
		"sqlite3", driver)
	m.Up()

	if err != nil {
		errorLog.Println(err)
	}
}
