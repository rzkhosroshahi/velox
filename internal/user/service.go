package user

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userStore *UserStore
	logger    *zap.Logger
}

func NewService(userStore *UserStore, logger *zap.Logger) *Service {
	return &Service{userStore: userStore, logger: logger}
}

func (s *Service) CreateUser(ctx context.Context, params CreateUserRequest) (User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(err.Error())
		return User{}, err
	}

	return s.userStore.CreateUserWithIdentity(ctx,
		&User{
			Name:  params.Name,
			Email: params.Email,
		},
		&UserIdentity{
			Password: string(hashed),
		},
	)
}
