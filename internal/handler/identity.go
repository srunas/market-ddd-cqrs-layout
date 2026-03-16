package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
	"github.com/srunas/market-ddd-cqrs-layout/internal/handler/middleware"
)

type IdentityHandler struct {
	svc service.Identity
}

func NewIdentityHandler(svc service.Identity) *IdentityHandler {
	return &IdentityHandler{svc: svc}
}

func (h *IdentityHandler) Register(w http.ResponseWriter, r *http.Request) {
	log := middleware.FromContext(r.Context())

	var req struct {
		Username string `json:"username"`
		Surname  string `json:"surname"`
		Email    string `json:"email"`
		Password string `json:"password"` //nolint:gosec // поле передаётся пользователем через JSON
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "некорректный запрос", http.StatusBadRequest)
		return
	}

	resp, err := h.svc.RegisterUser(r.Context(), service.RegisterUserRequest{
		Username: req.Username,
		Surname:  req.Surname,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Error("ошибка регистрации", "error", err)
		http.Error(w, "ошибка регистрации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"user_id": uuid.UUID(resp.UserID).String(),
	})
}

func (h *IdentityHandler) Login(w http.ResponseWriter, r *http.Request) {
	log := middleware.FromContext(r.Context())

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"` //nolint:gosec // поле передаётся пользователем через JSON
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "некорректный запрос", http.StatusBadRequest)
		return
	}

	resp, err := h.svc.LoginUser(r.Context(), service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.Error("ошибка входа", "error", err)
		http.Error(w, "неверные данные", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": resp.Token,
	})
}
