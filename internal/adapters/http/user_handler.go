package http

import (
	"net/http"

	"github.com/Longseeyou/DuraFlow/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserHandler struct {
	service *user.UserService
}

func NewUserHandler(service *user.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.createUserHandler)

	return r
}

func (h *UserHandler) ProtectedRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/me", h.getUserHandler)
	r.Put("/me", h.updateUserHandler)
	r.Delete("/me", h.deleteUserHandler)

	return r
}

func (h *UserHandler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var request user.CreateUserRequestDto
	if err := render.DecodeJSON(r.Body, &request); err != nil {
		renderError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.service.CreateUser(r.Context(), request)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, "failed to create user")
		return
	}

	renderJSON(w, r, http.StatusCreated, response)
}

func (h *UserHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDFromToken(w, r)
	if !ok {
		return
	}

	response, err := h.service.GetUser(r.Context(), user.UserRequestDto{ID: &id})
	if err != nil {
		renderError(w, r, http.StatusNotFound, "user not found")
		return
	}

	renderJSON(w, r, http.StatusOK, response)
}

func (h *UserHandler) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDFromToken(w, r)
	if !ok {
		return
	}

	var request user.UpdateUserRequestDto
	if err := render.DecodeJSON(r.Body, &request); err != nil {
		renderError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	request.ID = &id

	response, err := h.service.UpdateUser(r.Context(), request)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, "failed to update user")
		return
	}

	renderJSON(w, r, http.StatusOK, response)
}

func (h *UserHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDFromToken(w, r)
	if !ok {
		return
	}

	request := user.DeleteUserRequestDto{ID: &id}
	var (
		response user.UserResponseDto
		err      error
	)
	// soft delete first cases later
	response, err = h.service.SoftDeleteUser(r.Context(), request)

	if err != nil {
		renderError(w, r, http.StatusInternalServerError, "failed to delete user")
		return
	}

	renderJSON(w, r, http.StatusOK, response)
}

func renderJSON(w http.ResponseWriter, r *http.Request, statusCode int, response any) {
	render.Status(r, statusCode)
	render.JSON(w, r, response)
}

func renderError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	renderJSON(w, r, statusCode, map[string]string{"error": message})
}
