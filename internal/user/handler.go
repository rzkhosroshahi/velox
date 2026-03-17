package user

import (
	"encoding/json"
	"net/http"

	"github.com/rzkhosroshahi/velox/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	req := CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusInternalServerError, "invalid request body", nil)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusInternalServerError, "name, email and password are required", nil)
		return
	}
	user, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create user", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}
