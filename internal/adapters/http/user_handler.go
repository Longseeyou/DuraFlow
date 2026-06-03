package http

import (
	"net/http"

	"github.com/Longseeyou/DuraFlow/internal/user"
)

type UserHandler struct {
	service *user.UserService
}

func (h *UserHandler) createUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) updateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {

}
