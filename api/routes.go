package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rzkhosroshahi/velox/internal/token"
	"github.com/rzkhosroshahi/velox/internal/user"
)

func NewRouter(userHandler *user.Handler, tokenHandler *token.Handler) chi.Router {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", userHandler.Routes())

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(tokenHandler.Service))
			r.Mount("/auth", tokenHandler.Routes())
		})
	})

	r.Get("/health", HealthCheck)
	return r
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "status is available\n")
}
