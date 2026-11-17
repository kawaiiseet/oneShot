package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

type Server struct {
	repo PersonalDataRepository
}

func New(repo PersonalDataRepository) *Server {
	return &Server{repo: repo}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", s.handleRegister)
	mux.HandleFunc("POST /find", s.handleFind)
	return mux
}

// ---- /register ---------------------------------------------------------

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var p PersonalData
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if msg := validate(p); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	if err := s.repo.Create(r.Context(), &p); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			writeError(w, http.StatusConflict, "login already taken")
			return
		}
		log.Printf("register: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

// ---- /find -------------------------------------------------------------

func (s *Server) handleFind(w http.ResponseWriter, r *http.Request) {
	var req FindRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	country := req.Country
	interests := req.Interests
	login := req.Login

	if country == "" || len(interests) == 0 {
		writeError(w, http.StatusBadRequest, "country and interests required")
		return
	}

	match, err := s.repo.FindMatch(r.Context(), country, interests, login)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "no match found")
		return
	}
	if err != nil {
		log.Printf("find: match: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, match)
}
