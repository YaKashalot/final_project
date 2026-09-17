package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var database *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX idx_scheduler_date ON scheduler (date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return fmt.Errorf("failed to check db file: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	if install {
		if _, err := db.Exec(schema); err != nil {
			db.Close()
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}

	database = db
	return nil
}

func GetDB() *sql.DB {
	return database
}

func Close() error {
	if database != nil {
		return database.Close()
	}
	return nil
}
