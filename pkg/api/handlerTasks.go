package api

import (
	"net/http"

	"final/v1/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(r.Context(), 50, search)
	if err != nil {
		writeDBError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}
