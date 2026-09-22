package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type signInRequest struct {
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

const jwtSecret = "final_project_secret_key"

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req signInRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, signInResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	storedPassword := os.Getenv("TODO_PASSWORD")

	if req.Password != storedPassword {
		writeJSON(w, signInResponse{Error: "invalid password"})
		return
	}

	hash := passwordHash(storedPassword)

	claims := jwt.MapClaims{
		"hash": hash,
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		writeJSON(w, signInResponse{Error: "failed to sign token: " + err.Error()})
		return
	}

	writeJSON(w, signInResponse{Token: signed})
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", sum)
}
