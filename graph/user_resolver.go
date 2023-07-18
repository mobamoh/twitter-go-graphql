package graph

import (
	"context"
	twitter_go_graphql "github.com/mobamoh/twitter-go-graphql"
)

func (q *queryResolver) Me(ctx context.Context) (*User, error) {
	userIDFromContext, err := twitter_go_graphql.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, twitter_go_graphql.ErrUnauthenticated
	}

	return mapUser(twitter_go_graphql.User{
		ID: userIDFromContext,
	}), nil
}
