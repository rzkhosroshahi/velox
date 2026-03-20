package user

import (
	"context"

	"github.com/rzkhosroshahi/velox/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userStore *UserStore
}

func NewService(userStore *UserStore) *Service {
	return &Service{userStore: userStore}
}

func (s *Service) CreateUser(ctx context.Context, params CreateUserRequest) (User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error(err.Error())
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
