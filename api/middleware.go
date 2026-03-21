package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/rzkhosroshahi/velox/internal/token"
	"github.com/rzkhosroshahi/velox/pkg/logger"
	"github.com/rzkhosroshahi/velox/pkg/response"
	"go.uber.org/zap"
)

func AuthMiddleware(tokenService *token.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "missing token", nil)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := tokenService.Validate(r.Context(), tokenStr)
			if err != nil {
				logger.Log.Error("invalid token", zap.Error(err))
				response.Error(w, http.StatusUnauthorized, "invalid token", err)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			ctx = context.WithValue(ctx, "sessionID", claims.SessionID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
