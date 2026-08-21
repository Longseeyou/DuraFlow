package http

import (
	"errors"
	"net/http"

	"github.com/Longseeyou/DuraFlow/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type AuthHandler struct {
	service *auth.AuthService
}

func NewAuthHandler(service *auth.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", h.loginHandler)

	return r
}

func (h *AuthHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	var request auth.LoginRequestDto
	if err := render.DecodeJSON(r.Body, &request); err != nil {
		renderError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.service.Login(r.Context(), request)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			renderError(w, r, http.StatusUnauthorized, "invalid email or password")
			return
		}
		renderError(w, r, http.StatusInternalServerError, "failed to login")
		return
	}

	renderJSON(w, r, http.StatusOK, response)
}
