package api

import (
	"errors"
	"time"

	"final/v1/pkg/date"
	"final/v1/pkg/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return errors.New("invalid date format, expected 20060102")
	}

	var next string
	if task.Repeat != "" {
		next, err = date.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// прошедшую дату переносим на сегодня или по правилу повторения
	if isBeforeToday(t, now) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}

func isBeforeToday(value, today time.Time) bool {
	dy, dm, dd := value.Date()
	ny, nm, nd := today.Date()
	d := time.Date(dy, dm, dd, 0, 0, 0, 0, time.UTC)
	n := time.Date(ny, nm, nd, 0, 0, 0, 0, time.UTC)
	return d.Before(n)
}
