package http

import (
	"encoding/json"
	"net/http"

	"ancianato-backend/internal/infrastructure/auth"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct{}

func NewAuthHandler(r chi.Router) {
	handler := &AuthHandler{}
	r.Post("/login", handler.Login)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	token, err := auth.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Credenciales inválidas: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
