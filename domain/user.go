package domain

import (
	"context"
	twitter_go_graphql "github.com/mobamoh/twitter-go-graphql"
)

type UserService struct {
	UserRepo twitter_go_graphql.UserRepo
}

func NewUserService(repo twitter_go_graphql.UserRepo) *UserService {
	return &UserService{
		UserRepo: repo,
	}
}

func (u *UserService) GetByID(ctx context.Context, id string) (twitter_go_graphql.User, error) {
	return u.UserRepo.GetByID(ctx, id)
}
