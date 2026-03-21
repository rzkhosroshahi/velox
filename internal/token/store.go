package token

import (
	"github.com/golang-jwt/jwt/v5"
)

const (
	ScopeAuth = "authentication"
)

type Claims struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	Scope     string `json:"scope"`
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
	RefreshToken string `json:"refreshToken"`
}

type GenerateParams struct {
	UserID    string
	UserAgent string
	IPAddress string
}
