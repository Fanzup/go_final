package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func decodeJSON(w http.ResponseWriter, r *http.Request, value any) error {
	// не читаем больше 1 МБ
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	// после JSON допускаются только пробелы
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			return err
		}
		return errors.New("expected a single JSON object")
	}
	return nil
}
