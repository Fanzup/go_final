package api

import (
	"net/http"

	"final/v1/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "task ID is required"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeDBError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}
