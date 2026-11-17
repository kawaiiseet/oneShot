package server

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func validate(p PersonalData) string {
	switch {
	case p.Name == "":
		return "name is required"
	case p.Login == "":
		return "login is required"
	case p.Age < 18 || p.Age > 150:
		return "age out of range"
	case p.Country == "":
		return "country is required"
	case p.Password == "":
		return "password is required"
	}
	return ""
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func passwordHashing(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
