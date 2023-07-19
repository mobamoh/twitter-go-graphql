package graph

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	twitter_go_graphql "github.com/mobamoh/twitter-go-graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
)

//go:generate go get -d github.com/99designs/gqlgen
//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	AuthService  twitter_go_graphql.AuthService
	TweetService twitter_go_graphql.TweetService
	UserService  twitter_go_graphql.UserService
}

type tweetResolver struct {
	*Resolver
}

func (r *Resolver) Tweet() TweetResolver {
	return &tweetResolver{r}
}

type queryResolver struct {
	*Resolver
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct {
	*Resolver
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func buildBadRequestError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusBadRequest,
		},
	}
}

func buildUnauthenticatedError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusUnauthorized,
		},
	}
}

func buildForbiddenError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusForbidden,
		},
	}
}

func buildError(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, twitter_go_graphql.ErrBadCredentials):
		return buildUnauthenticatedError(ctx, err)
	case errors.Is(err, twitter_go_graphql.ErrValidation) || errors.Is(err, twitter_go_graphql.ErrInvalidUUID):
		return buildBadRequestError(ctx, err)
	case errors.Is(err, twitter_go_graphql.ErrForbidden):
		return buildForbiddenError(ctx, err)
	default:
		return fmt.Errorf("%+v", err)
	}
}
