package api

import "net/http"

func Init(mux *http.ServeMux) error {
	config, err := loadAuthConfig()
	if err != nil {
		return err
	}
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", config.auth(taskHandler))
	mux.HandleFunc("/api/tasks", config.auth(tasksHandler))
	mux.HandleFunc("/api/task/done", config.auth(doneTaskHandler))
	mux.HandleFunc("/api/signin", config.signinHandler)
	return nil
}
