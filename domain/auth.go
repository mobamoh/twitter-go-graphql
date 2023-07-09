package domain

import (
	"context"
	"errors"
	twitter "github.com/mobamoh/twitter-go-graphql"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo twitter.UserRepo
}

func NewAuthService(repo twitter.UserRepo) *AuthService {
	return &AuthService{
		userRepo: repo,
	}
}

func (svc *AuthService) Register(ctx context.Context, input twitter.RegisterInput) (twitter.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitter.AuthResponse{}, err
	}

	if _, err := svc.userRepo.GetByUsername(ctx, input.UserName); !errors.Is(err, twitter.ErrNotFound) {
		return twitter.AuthResponse{}, twitter.ErrUsernameTaken
	}

	if _, err := svc.userRepo.GetByEmail(ctx, input.Email); !errors.Is(err, twitter.ErrNotFound) {
		return twitter.AuthResponse{}, twitter.ErrEmailTaken
	}

	user := twitter.User{
		Username: input.UserName,
		Email:    input.Email,
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return twitter.AuthResponse{}, twitter.ErrServer
	}
	user.Password = string(hashPwd)
	user, err = svc.userRepo.Create(ctx, user)
	if err != nil {
		return twitter.AuthResponse{}, twitter.ErrServer
	}
	return twitter.AuthResponse{
		AccessToken: "token",
		User:        user,
	}, nil
}
