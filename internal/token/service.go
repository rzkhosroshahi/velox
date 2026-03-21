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

func (s *Service) sign(userId string, sessionID string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID:    userId,
		SessionID: sessionID,
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
