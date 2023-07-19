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

func TestTweetService_Delete(t *testing.T) {
	t.Run("not auth user cannot delete a tweet", func(t *testing.T) {
		ctx := context.Background()
		err := tweetService.Delete(ctx, faker.UUID())
		require.ErrorIs(t, err, twitter_go_graphql.ErrUnauthenticated)
	})
	t.Run("cannot delete tweet if not owner", func(t *testing.T) {
		ctx := context.Background()
		test_helper.TeardownDB(ctx, t, db)

		otherUser := test_helper.CreateUser(ctx, t, userRepo)
		currentUser := test_helper.CreateUser(ctx, t, userRepo)

		existingTweet := test_helper.CreateTweet(ctx, t, tweetRepo, otherUser.ID)
		ctx = test_helper.LoginUser(ctx, t, currentUser)
		_, err := tweetRepo.GetByID(ctx, existingTweet.ID)
		require.NoError(t, err)

		err = tweetService.Delete(ctx, existingTweet.ID)
		require.ErrorIs(t, err, twitter_go_graphql.ErrForbidden)

		_, err = tweetRepo.GetByID(ctx, existingTweet.ID)
		require.NoError(t, err)
	})
	t.Run("can delete tweet", func(t *testing.T) {
		ctx := context.Background()
		test_helper.TeardownDB(ctx, t, db)

		currentUser := test_helper.CreateUser(ctx, t, userRepo)

		existingTweet := test_helper.CreateTweet(ctx, t, tweetRepo, currentUser.ID)
		ctx = test_helper.LoginUser(ctx, t, currentUser)

		_, err := tweetRepo.GetByID(ctx, existingTweet.ID)
		require.NoError(t, err)

		err = tweetService.Delete(ctx, existingTweet.ID)
		require.NoError(t, err)

		_, err = tweetRepo.GetByID(ctx, existingTweet.ID)
		require.ErrorIs(t, err, twitter_go_graphql.ErrNotFound)
	})
	t.Run("return error invalid uuid ", func(t *testing.T) {
		ctx := context.Background()
		test_helper.TeardownDB(ctx, t, db)

		currentUser := test_helper.CreateUser(ctx, t, userRepo)
		ctx = test_helper.LoginUser(ctx, t, currentUser)

		err := tweetService.Delete(ctx, "123")
		require.ErrorIs(t, err, twitter_go_graphql.ErrInvalidUUID)

	})
}

func TestTweetService_CreateReply(t *testing.T) {
	t.Run("not auth user cannot reply to a tweet", func(t *testing.T) {
		ctx := context.Background()

		input := twitter_go_graphql.CreateTweetInput{
			Body: faker.RandStringRunes(20),
		}
		_, err := tweetService.CreateReply(ctx, faker.UUID(), input)
		require.ErrorIs(t, err, twitter_go_graphql.ErrUnauthenticated)
	})
	t.Run("cannot reply to non existing tweet", func(t *testing.T) {
		ctx := context.Background()
		test_helper.TeardownDB(ctx, t, db)

		currentUser := test_helper.CreateUser(ctx, t, userRepo)
		ctx = test_helper.LoginUser(ctx, t, currentUser)

		input := twitter_go_graphql.CreateTweetInput{
			Body: faker.RandStringRunes(20),
		}
		_, err := tweetService.CreateReply(ctx, faker.UUID(), input)
		require.ErrorIs(t, err, twitter_go_graphql.ErrNotFound)
	})
	t.Run("can reply to  tweet", func(t *testing.T) {
		ctx := context.Background()
		test_helper.TeardownDB(ctx, t, db)

		currentUser := test_helper.CreateUser(ctx, t, userRepo)
		ctx = test_helper.LoginUser(ctx, t, currentUser)

		input := twitter_go_graphql.CreateTweetInput{
			Body: faker.RandStringRunes(20),
		}

		existingTweet := test_helper.CreateTweet(ctx, t, tweetRepo, currentUser.ID)

		reply, err := tweetService.CreateReply(ctx, existingTweet.ID, input)

		require.NoError(t, err)
		require.NotEmpty(t, reply.ID)
		require.Equal(t, input.Body, reply.Body)
		require.Equal(t, currentUser.ID, reply.UserID)
		require.Equal(t, existingTweet.ID, *reply.ParentID)
		require.NotEmpty(t, reply.CreatedAt)

	})
}
