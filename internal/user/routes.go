package user

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	return r
}
