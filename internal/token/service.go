package token

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	accessTTL  = 15 * time.Minute
	refreshTTL = 1 * time.Hour
)

type Service struct {
	redis     *redis.Client
	secretKey string
}

func NewService(redis *redis.Client, secretKey string) *Service {
	return &Service{
		redis:     redis,
		secretKey: secretKey,
	}
}

func (s *Service) Generate(ctx context.Context, userID, userAgent, ipAddress string) (TokenPair, error) {
	sessionID := uuid.New().String()

	accessToken, err := s.sign(userID, sessionID, accessTTL)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := s.sign(userID, sessionID, refreshTTL)
	if err != nil {
		return TokenPair{}, err
	}
	err = s.redis.Set(ctx, accessKey(sessionID), accessToken, accessTTL).Err()
	if err != nil {
		return TokenPair{}, err
	}

	err = s.redis.Set(ctx, refreshKey(sessionID), refreshToken, refreshTTL).Err()
	if err != nil {
		return TokenPair{}, err
	}

	s.redis.HSet(ctx, sessionMetaKey(sessionID), map[string]any{
		"sessionId": sessionID,
		"userId":    userID,
		"userAgent": userAgent,
		"ipAddress": ipAddress,
		"scope":     ScopeAuth,
		"createdAt": time.Now().Format(time.RFC3339),
	})
	s.redis.Expire(ctx, sessionsKey(userID), refreshTTL)
	s.redis.SAdd(ctx, sessionsKey(userID), sessionID)
	s.redis.Expire(ctx, sessionsKey(userID), refreshTTL)

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (s *Service) Refresh(ctx context.Context, refreshToken, userAgent, ipAddress string) (TokenPair, error) {
	claims, err := s.parse(refreshToken)
	if err != nil {
		return TokenPair{}, err
	}

	stored, err := s.redis.Get(ctx, refreshKey(claims.SessionID)).Result()
	if err != nil || stored != refreshToken {
		return TokenPair{}, fmt.Errorf("refresh token expired or revoked")
	}

	if err := s.revokeSession(ctx, claims.UserID, claims.SessionID); err != nil {
		return TokenPair{}, err
	}

	return s.Generate(ctx, claims.UserID, userAgent, ipAddress)
}
func (s *Service) Revoke(ctx context.Context, userID, sessionID string) error {
	s.revokeSession(ctx, userID, sessionID)
	return nil
}
func (s *Service) revokeSession(ctx context.Context, userID, sessionID string) error {
	s.redis.Del(ctx, accessKey(sessionID))
	s.redis.Del(ctx, refreshKey(sessionID))
	s.redis.Del(ctx, sessionMetaKey(sessionID))
	s.redis.SRem(ctx, sessionsKey(userID), sessionID)
	return nil
}
func (s *Service) Validate(ctx context.Context, accessToken string) (*Claims, error) {
	claims, err := s.parse(accessToken)
	if err != nil {
		return nil, err
	}

	stored, err := s.redis.Get(ctx, accessKey(claims.SessionID)).Result()
	if err != nil || stored != accessToken {
		return nil, fmt.Errorf("token expired or revoked")
	}

	return claims, nil
}
func (s *Service) sign(userId string, sessionID string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID:    userId,
		SessionID: sessionID,
		Scope:     ScopeAuth,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.secretKey))
}
func (s *Service) parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token.Claims.(*Claims), nil
}

func accessKey(sessionID string) string      { return "access:" + sessionID }
func refreshKey(sessionID string) string     { return "refresh:" + sessionID }
func sessionsKey(userID string) string       { return "sessions:" + userID }
func sessionMetaKey(sessionID string) string { return "session:meta:" + sessionID }
