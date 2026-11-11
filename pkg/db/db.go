package db

import (
	"database/sql"

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
	conn, err := sql.Open("sqlite", fileDB)
	if err != nil {
		return err
	}
	// одно соединение: операции с БД идут по очереди
	conn.SetMaxOpenConns(1)
	// создаём таблицу и индекс, если их ещё нет
	if _, err := conn.Exec(schema); err != nil {
		_ = conn.Close()
		return err
	}
	db = conn
	return nil
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
