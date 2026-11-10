package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"final/v1/pkg/db"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	buf, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(buf)
}

func methodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
}

func writeDBError(w http.ResponseWriter, err error) {
	if errors.Is(err, db.ErrTaskNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	log.Printf("database operation failed: %v", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to perform task operation"})
}
