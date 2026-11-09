package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"final/v1/pkg/date"
	"final/v1/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "task title is required"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

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

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}

func afterNow(date, now time.Time) bool {
	dy, dm, dd := date.Date()
	ny, nm, nd := now.Date()
	d := time.Date(dy, dm, dd, 0, 0, 0, 0, time.UTC)
	n := time.Date(ny, nm, nd, 0, 0, 0, 0, time.UTC)
	return d.After(n)
}
