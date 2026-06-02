package http

import (
	"net/http"

	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/go-chi/chi/v5"
)

type RouterConfig struct {
	WorkflowService *workflow.WorkflowService
}

func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {

	})
	return r
}
