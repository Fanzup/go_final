package api

import (
	"net/http"
	"strconv"

	"final/v1/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := decodeJSON(w, r, &task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "task title is required"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeDBError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}
