package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	Name      string    `db:"name"       json:"name"`
	Email     string    `db:"email"      json:"email"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type UserIdentity struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	UserID    uuid.UUID `db:"user_id"    json:"userId"`
	Provider  string    `db:"provider"   json:"provider"`
	Password  string    `db:"password"   json:"-"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
	Password string `json:"password"`
}

type LoginParams struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserAgent string `json:"userAgent"`
	IPAddress string `json:"IPAddress"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

func (us *UserStore) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	var user User
	err := us.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	return user, err
}

func (us *UserStore) GetUserIdentityByUserID(ctx context.Context, id uuid.UUID) (UserIdentity, error) {
	var userIdentity UserIdentity
	err := us.db.GetContext(ctx, &userIdentity, "SELECT * FROM user_identities WHERE user_id = $1", id)
	return userIdentity, err
}

func (us *UserStore) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := us.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	return user, err
}

func (us *UserStore) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := us.db.QueryRowxContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email,
	).Scan(&exists)
	return exists, err
}

func (us *UserStore) CreateUserWithIdentity(ctx context.Context, user *User, identity *UserIdentity) (User, error) {
	taken, err := us.IsEmailTaken(ctx, user.Email)
	if err != nil {
		return User{}, err
	}
	if taken {
		return User{}, errors.New("email already taken")
	}

	tx, err := us.db.BeginTxx(ctx, nil)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback()

	userQuery := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at
	`
	createdUser := User{}
	err = tx.QueryRowxContext(ctx, userQuery, user.Name, user.Email).StructScan(&createdUser)
	if err != nil {
		return User{}, err
	}

	identityQuery := `
		INSERT INTO user_identities (user_id, password)
		VALUES ($1, $2)
	`
	_, err = tx.ExecContext(ctx, identityQuery, createdUser.ID, identity.Password)
	if err != nil {
		return User{}, err
	}

	if err := tx.Commit(); err != nil {
		return User{}, err
	}

	return createdUser, nil
}

func (us *UserStore) GetIdentityByEmail(ctx context.Context, email string) (UserIdentity, error) {
	user, err := us.GetUserByEmail(ctx, email)
	if err != nil {
		return UserIdentity{}, err
	}
	identity, err := us.GetUserIdentityByUserID(ctx, user.ID)
	if err != nil {
		return UserIdentity{}, err
	}
	return identity, err
}
