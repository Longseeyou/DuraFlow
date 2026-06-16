package http

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

func userIDFromToken(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		renderError(w, r, http.StatusUnauthorized, "invalid token")
		return uuid.Nil, false
	}

	rawID, ok := claims["user_id"].(string)
	if !ok || rawID == "" {
		rawID, ok = claims["sub"].(string)
	}
	if !ok || rawID == "" {
		renderError(w, r, http.StatusUnauthorized, "missing user id claim")
		return uuid.Nil, false
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		renderError(w, r, http.StatusUnauthorized, "invalid user id claim")
		return uuid.Nil, false
	}

	return id, true
}
