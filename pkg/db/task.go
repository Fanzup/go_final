package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"final/v1/pkg/date"
)

var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?,?,?,?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	var t Task
	query := `SELECT CAST(id AS TEXT), date, title, comment, repeat
	          FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &t, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	if limit < 1 {
		return []*Task{}, nil
	}
	query := `SELECT CAST(id AS TEXT), date, title, comment, repeat FROM scheduler`
	args := []any{}
	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			query += ` WHERE date = ?`
			args = append(args, t.Format("20060102"))
		} else {
			query += ` WHERE title LIKE ? OR comment LIKE ?`
			pattern := "%" + search + "%"
			args = append(args, pattern, pattern)
		}
	}
	query += ` ORDER BY date, id LIMIT ?`
	args = append(args, limit)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0, limit)
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err()
}

// читаем и завершаем задачу в одной транзакции
func CompleteTask(ctx context.Context, id string, now time.Time) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var task Task
	err = tx.QueryRowContext(ctx, `SELECT date, repeat FROM scheduler WHERE id = ?`, id).Scan(&task.Date, &task.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrTaskNotFound
	}
	if err != nil {
		return err
	}
	// одноразовую удаляем, повторяющуюся переносим
	if task.Repeat == "" {
		_, err = tx.ExecContext(ctx, `DELETE FROM scheduler WHERE id = ?`, id)
	} else {
		next, dateErr := date.NextDate(now, task.Date, task.Repeat)
		if dateErr != nil {
			return dateErr
		}
		_, err = tx.ExecContext(ctx, `UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}
