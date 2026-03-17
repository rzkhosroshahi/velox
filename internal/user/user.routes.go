package user

import "github.com/go-chi/chi/v5"

func UsersRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", GetUsers)

	return r
}
