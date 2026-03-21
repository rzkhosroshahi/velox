package user

import (
	"context"
	"errors"

	"github.com/rzkhosroshahi/velox/internal/token"
	"github.com/rzkhosroshahi/velox/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store      *UserStore
	tokenStore *token.Service
}

func NewService(userStore *UserStore) *Service {
	return &Service{store: userStore}
}

func (s *Service) CreateUser(ctx context.Context, params CreateUserRequest) (User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error(err.Error())
		return User{}, err
	}

	return s.store.CreateUserWithIdentity(ctx,
		&User{
			Name:  params.Name,
			Email: params.Email,
		},
		&UserIdentity{
			Password: string(hashed),
		},
	)
}

func (s *Service) LoginUser(ctx context.Context, params LoginParams) (User, error) {
	identity, err := s.store.GetIdentityByEmail(ctx, params.Email)
	if err != nil {
		return User{}, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(identity.Password), []byte(params.Password)); err != nil {
		return User{}, errors.New("invalid credentials")
	}

	return s.store.GetUserByID(ctx, identity.UserID)
}
