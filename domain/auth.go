package domain

import (
	"context"
	"errors"
	"fmt"
	twitter "github.com/mobamoh/twitter-go-graphql"
	"golang.org/x/crypto/bcrypt"
)

var PasswordHashCost = bcrypt.DefaultCost

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

	if _, err := svc.userRepo.GetByUsername(ctx, input.Username); !errors.Is(err, twitter.ErrNotFound) {
		return twitter.AuthResponse{}, twitter.ErrUsernameTaken
	}

	if _, err := svc.userRepo.GetByEmail(ctx, input.Email); !errors.Is(err, twitter.ErrNotFound) {
		return twitter.AuthResponse{}, twitter.ErrEmailTaken
	}

	user := twitter.User{
		Username: input.Username,
		Email:    input.Email,
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), PasswordHashCost)
	if err != nil {
		return twitter.AuthResponse{}, twitter.ErrServer
	}
	user.Password = string(hashPwd)
	user, err = svc.userRepo.Create(ctx, user)
	if err != nil {
		return twitter.AuthResponse{}, fmt.Errorf("%w: %v", twitter.ErrServer, err)
	}
	return twitter.AuthResponse{
		AccessToken: "token",
		User:        user,
	}, nil
}

func (svc *AuthService) Login(ctx context.Context, input twitter.LoginInput) (twitter.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitter.AuthResponse{}, err
	}

	user, err := svc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, twitter.ErrNotFound):
			return twitter.AuthResponse{}, twitter.ErrBadCredentials
		default:
			return twitter.AuthResponse{}, err
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return twitter.AuthResponse{}, twitter.ErrBadCredentials
	}
	return twitter.AuthResponse{
		AccessToken: "token",
		User:        user,
	}, nil
}
