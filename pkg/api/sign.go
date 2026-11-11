package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (config authConfig) signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	if config.password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "password authentication is disabled"})
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	supplied := sha256.Sum256([]byte(req.Password))
	expected := sha256.Sum256([]byte(config.password))
	if subtle.ConstantTimeCompare(supplied[:], expected[:]) != 1 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": config.passwordHash(),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})
	signed, err := token.SignedString(config.secret)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create token"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": signed})
}
