package http

import (
	"net/http"

	"github.com/Longseeyou/DuraFlow/internal/auth"
	"github.com/Longseeyou/DuraFlow/internal/user"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type RouterConfig struct {
	WorkflowService *workflow.WorkflowService
	UserService     *user.UserService
	AuthService     *auth.AuthService
}

func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)

	r.Get("/healthz", Healthz)
	r.Get("/readyz", Readyz)

	if cfg.AuthService != nil {
		r.Mount("/auth", NewAuthHandler(cfg.AuthService).Routes())

		if cfg.UserService != nil {
			userHandler := NewUserHandler(cfg.UserService)
			r.Post("/users", userHandler.createUserHandler)
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(cfg.AuthService.TokenAuth))
				r.Use(jwtauth.Authenticator(cfg.AuthService.TokenAuth))
				r.Mount("/users", userHandler.ProtectedRoutes())
			})
		}
	} else if cfg.UserService != nil {
		r.Mount("/users", NewUserHandler(cfg.UserService).Routes())
	}

	return r
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{
		"status": "ok",
	})
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{
		"status": "ready",
	})
}
