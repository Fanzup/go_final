package api

import (
	"net/http"

	"final/v1/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "task ID is required"})
		return
	}

	task, err := db.GetTask(r.Context(), id)
	if err != nil {
		writeDBError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}
