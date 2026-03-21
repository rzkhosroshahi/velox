package user

import (
	"encoding/json"
	"net/http"

	"github.com/rzkhosroshahi/velox/internal/token"
	"github.com/rzkhosroshahi/velox/pkg/logger"
	"github.com/rzkhosroshahi/velox/pkg/response"
	"go.uber.org/zap"
)

type Handler struct {
	service      *Service
	tokenService *token.Service
}

func NewHandler(service *Service, tokenService *token.Service) *Handler {
	return &Handler{service: service, tokenService: tokenService}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	req := CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("json.NewDecoder invalid request body")
		response.Error(w, http.StatusInternalServerError, "invalid request body", nil)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		logger.Log.Error("validation user request request body")
		response.Error(w, http.StatusInternalServerError, "name, email and password are required", nil)
		return
	}
	user, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		logger.Log.Error("CreateUser user", zap.Error(err))
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	user, err := h.service.LoginUser(r.Context(), LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials", err)
		return
	}

	pair, err := h.tokenService.Generate(r.Context(), user.ID.String(), r.Header.Get("User-Agent"), r.RemoteAddr)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to generate token", err)
		return
	}

	response.JSON(w, http.StatusOK, pair)
}
