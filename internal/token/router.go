package token

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/logout", h.Logout)
	r.Post("/refresh", h.Refresh)

	return r
}
