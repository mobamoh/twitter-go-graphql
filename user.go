package twitter_go_graphql

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUsernameTaken = errors.New("username taken")
	ErrEmailTaken    = errors.New("email taken")
)

type UserService interface {
	GetByID(ctx context.Context, id string) (User, error)
}
type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	GetByIDs(ctx context.Context, id []string) ([]User, error)
}
type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
