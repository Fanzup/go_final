package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type authConfig struct {
	password string
	secret   []byte
}

func loadAuthConfig() (authConfig, error) {
	config := authConfig{password: os.Getenv("TODO_PASSWORD"), secret: []byte(os.Getenv("TODO_JWT_SECRET"))}
	if config.password != "" && len(config.secret) < 32 {
		return authConfig{}, fmt.Errorf("TODO_JWT_SECRET must contain at least 32 bytes when TODO_PASSWORD is set")
	}
	return config, nil
}

func (config authConfig) passwordHash() string {
	// смена пароля делает старые токены недействительными
	mac := hmac.New(sha256.New, config.secret)
	mac.Write([]byte("scheduler-password-version:"))
	mac.Write([]byte(config.password))
	return hex.EncodeToString(mac.Sum(nil))
}

func (config authConfig) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.password == "" {
			next(w, r)
			return
		}
		unauthorized := func() {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			unauthorized()
			return
		}
		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
			return config.secret, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
		if err != nil || !token.Valid {
			unauthorized()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			unauthorized()
			return
		}
		hash, _ := claims["hash"].(string)
		if subtle.ConstantTimeCompare([]byte(hash), []byte(config.passwordHash())) != 1 {
			unauthorized()
			return
		}
		next(w, r)
	}
}
