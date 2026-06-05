package http

import (
	"net/http"

	"github.com/Longseeyou/DuraFlow/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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
	r.Get("/{id}", h.getUserHandler)
	r.Put("/{id}", h.updateUserHandler)
	r.Delete("/{id}", h.deleteUserHandler)

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
	id, ok := userIDFromRequest(w, r)
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
	id, ok := userIDFromRequest(w, r)
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
	id, ok := userIDFromRequest(w, r)
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

func userIDFromRequest(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	rawID := chi.URLParam(r, "id")
	if rawID == "" {
		rawID = chi.URLParam(r, "userID")
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, "invalid user id")
		return uuid.Nil, false
	}

	return id, true
}

func renderJSON(w http.ResponseWriter, r *http.Request, statusCode int, response any) {
	render.Status(r, statusCode)
	render.JSON(w, r, response)
}

func renderError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	renderJSON(w, r, statusCode, map[string]string{"error": message})
}
