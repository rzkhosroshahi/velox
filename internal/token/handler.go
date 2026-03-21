package token

import (
	"encoding/json"
	"net/http"

	"github.com/rzkhosroshahi/velox/pkg/logger"
	"github.com/rzkhosroshahi/velox/pkg/response"
	"go.uber.org/zap"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	sessionID := r.Context().Value("sessionID").(string)

	if err := h.Service.Revoke(r.Context(), userID, sessionID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to logout", err)
		return
	}

	response.JSON(w, http.StatusOK, "logged out successfully")
}
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	req := RefreshRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("failed to decode refresh request", zap.Error(err))
		response.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refresh token is required", nil)
		return
	}

	pair, err := h.Service.Refresh(r.Context(), req.RefreshToken, r.Header.Get("User-Agent"), r.RemoteAddr)
	if err != nil {
		logger.Log.Error("failed to refresh token", zap.Error(err))
		response.Error(w, http.StatusUnauthorized, "invalid or expired refresh token", nil)
		return
	}

	response.JSON(w, http.StatusOK, pair)
}
