//go:build integration
// +build integration

package domain

import (
	"context"
	twitter_go_graphql "github.com/mobamoh/twitter-go-graphql"
	"github.com/mobamoh/twitter-go-graphql/faker"
	"github.com/mobamoh/twitter-go-graphql/test_helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTweetService_Create(t *testing.T) {
	t.Run("unauthenticated user cannot create tweet", func(t *testing.T) {
		ctx := context.Background()
		_, err := tweetService.Create(ctx, twitter_go_graphql.CreateTweetInput{
			Body: "hello",
		})
		require.ErrorIs(t, err, twitter_go_graphql.ErrUnauthenticated)
	})
	t.Run("can create tweet", func(t *testing.T) {
		ctx := context.Background()
		test_helper.TeardownDB(ctx, t, db)
		user := test_helper.CreateUser(ctx, t, userRepo)
		ctx = test_helper.LoginUser(ctx, t, user)

		input := twitter_go_graphql.CreateTweetInput{
			Body: faker.RandStringRunes(200),
		}
		tweet, err := tweetService.Create(ctx, input)
		require.NoError(t, err)
		require.NotEmpty(t, tweet.ID)
		require.Equal(t, input.Body, tweet.Body)
		require.Equal(t, user.ID, tweet.UserID)
		require.NotEmpty(t, tweet.CreatedAt)
	})
}

func TestTweetService_GetByID(t *testing.T) {
	t.Run("can get tweet by id", func(t *testing.T) {
		ctx := context.Background()
		test_helper.TeardownDB(ctx, t, db)
		user := test_helper.CreateUser(ctx, t, userRepo)
		existingTweet := test_helper.CreateTweet(ctx, t, tweetRepo, user.ID)
		tweet, err := tweetService.GetByID(ctx, existingTweet.ID)
		require.NoError(t, err)
		require.Equal(t, existingTweet.ID, tweet.ID)
		require.Equal(t, existingTweet.Body, tweet.Body)
	})
	t.Run("return error not found", func(t *testing.T) {
		ctx := context.Background()
		test_helper.TeardownDB(ctx, t, db)

		_, err := tweetService.GetByID(ctx, faker.UUID())
		require.ErrorIs(t, err, twitter_go_graphql.ErrNotFound)

	})
	t.Run("return error invalid uuid ", func(t *testing.T) {
		ctx := context.Background()
		test_helper.TeardownDB(ctx, t, db)

		_, err := tweetService.GetByID(ctx, "123")
		require.ErrorIs(t, err, twitter_go_graphql.ErrInvalidUUID)

	})
}

func TestTweetService_All(t *testing.T) {
	t.Run("can get all tweet", func(t *testing.T) {
		ctx := context.Background()
		test_helper.TeardownDB(ctx, t, db)
		user := test_helper.CreateUser(ctx, t, userRepo)

		test_helper.CreateTweet(ctx, t, tweetRepo, user.ID)
		test_helper.CreateTweet(ctx, t, tweetRepo, user.ID)
		test_helper.CreateTweet(ctx, t, tweetRepo, user.ID)

		tweets, err := tweetService.All(ctx)
		require.NoError(t, err)
		require.Len(t, tweets, 3)
	})
}
