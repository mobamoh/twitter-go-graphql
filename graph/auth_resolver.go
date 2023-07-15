package graph

import (
	"context"
	"errors"
	twitter_go_graphql "github.com/mobamoh/twitter-go-graphql"
)

func mapUser(user twitter_go_graphql.User) *User {
	return &User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}
func mapAuthResponse(response twitter_go_graphql.AuthResponse) *AuthResponse {
	return &AuthResponse{
		AccessToken: response.AccessToken,
		User:        mapUser(response.User),
	}
}

func (m mutationResolver) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	res, err := m.AuthService.Register(ctx, twitter_go_graphql.RegisterInput{
		Email:           input.Email,
		Username:        input.Username,
		Password:        input.Password,
		ConfirmPassword: input.ConfirmPassword,
	})

	if err != nil {
		switch {
		case errors.Is(err, twitter_go_graphql.ErrValidation) ||
			errors.Is(err, twitter_go_graphql.ErrEmailTaken) ||
			errors.Is(err, twitter_go_graphql.ErrUsernameTaken):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}
	return mapAuthResponse(res), nil
}

func (m mutationResolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	res, err := m.AuthService.Login(ctx, twitter_go_graphql.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, twitter_go_graphql.ErrValidation) ||
			errors.Is(err, twitter_go_graphql.ErrBadCredentials):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}
	return mapAuthResponse(res), nil
}
