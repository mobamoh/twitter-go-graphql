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

func TestIntegrationAuthService_Register(t *testing.T) {
	validInput := twitter_go_graphql.RegisterInput{
		Username:        faker.Username(),
		Email:           faker.Email(),
		Password:        "password",
		ConfirmPassword: "password",
	}
	t.Run("can register", func(t *testing.T) {
		ctx := context.Background()
		defer test_helper.TeardownDB(ctx, t, db)
		res, err := authService.Register(ctx, validInput)

		require.NoError(t, err)
		require.NotEmpty(t, res.User.ID)
		require.Equal(t, validInput.Email, res.User.Email)
		require.Equal(t, validInput.Username, res.User.Username)
		require.NotEqual(t, validInput.Password, res.User.Password)
	})
	t.Run("existing username", func(t *testing.T) {
		ctx := context.Background()
		defer test_helper.TeardownDB(ctx, t, db)

		_, err := authService.Register(ctx, validInput)
		require.NoError(t, err)

		_, err = authService.Register(ctx, twitter_go_graphql.RegisterInput{
			Username:        validInput.Username,
			Email:           faker.Email(),
			Password:        "password",
			ConfirmPassword: "password",
		})
		require.ErrorIs(t, err, twitter_go_graphql.ErrUsernameTaken)
	})
	t.Run("existing email", func(t *testing.T) {
		ctx := context.Background()
		defer test_helper.TeardownDB(ctx, t, db)

		_, err := authService.Register(ctx, validInput)
		require.NoError(t, err)

		_, err = authService.Register(ctx, twitter_go_graphql.RegisterInput{
			Username:        faker.Username(),
			Email:           validInput.Email,
			Password:        "password",
			ConfirmPassword: "password",
		})
		require.ErrorIs(t, err, twitter_go_graphql.ErrEmailTaken)
	})
}
