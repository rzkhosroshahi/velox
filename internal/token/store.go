package token

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `son:"userId"`
	SessionID string ` json:"sessionId"`
	jwt.RegisteredClaims
}

type Session struct {
	SessionID string `json:"sessionId"`
	UserID    string `json:"userId"`
	UserAgent string `json:"userAgent"`
	IPAddress string `json:"ipAddress"`
	CreatedAt string `json:"createdAt"`
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type GenerateParams struct {
	UserID    string
	UserAgent string
	IPAddress string
}
