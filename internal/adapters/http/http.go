package http

import (
	"net/http"

	"github.com/go-chi/render"
)

func renderJSON(w http.ResponseWriter, r *http.Request, statusCode int, response any) {
	render.Status(r, statusCode)
	render.JSON(w, r, response)
}

func renderError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	renderJSON(w, r, statusCode, map[string]string{"error": message})
}
