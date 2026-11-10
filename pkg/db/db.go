package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR NOT NULL DEFAULT "",
	comment TEXT,
	repeat VARCHAR NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

var db *sql.DB

func InitDB(fileDB string) error {
	_, err := os.Stat(fileDB)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	install := os.IsNotExist(err)

	db, err = sql.Open("sqlite", fileDB)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return err
	}

	if install {
		if _, err = db.Exec(schema); err != nil {
			_ = db.Close()
			return err
		}
		log.Println("DB is init!")
	}

	return nil

}
