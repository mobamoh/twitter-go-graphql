package test_helper

import (
	"context"
	twitter_go_graphql "github.com/mobamoh/twitter-go-graphql"
	"github.com/mobamoh/twitter-go-graphql/faker"
	"github.com/mobamoh/twitter-go-graphql/postgres"
	"github.com/stretchr/testify/require"
	"testing"
)

func TeardownDB(ctx context.Context, t *testing.T, db *postgres.DB) {
	t.Helper()
	err := db.Truncate(ctx)
	require.NoError(t, err)
}

func CreateUser(ctx context.Context, t *testing.T, repo twitter_go_graphql.UserRepo) twitter_go_graphql.User {
	t.Helper()
	user, err := repo.Create(ctx, twitter_go_graphql.User{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.EncryptedPassword,
	})
	require.NoError(t, err)
	return user
}

func LoginUser(ctx context.Context, t *testing.T, user twitter_go_graphql.User) context.Context {
	t.Helper()
	return twitter_go_graphql.PutUserIDIntoContext(ctx, user.ID)
}

func CreateTweet(ctx context.Context, t *testing.T, repo twitter_go_graphql.TweetRepo, userId string) twitter_go_graphql.Tweet {
	t.Helper()
	tweet, err := repo.Create(ctx, twitter_go_graphql.Tweet{
		Body:   faker.RandStringRunes(33),
		UserID: userId,
	})
	require.NoError(t, err)
	return tweet
}
