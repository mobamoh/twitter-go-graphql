package twitter_go_graphql

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUserIDFromContext(t *testing.T) {
	t.Run("get user id from context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextAuthIDKey, "123")

		userIDFromContext, err := GetUserIDFromContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "123", userIDFromContext)
	})

	t.Run("return error if missing id", func(t *testing.T) {
		ctx := context.Background()
		_, err := GetUserIDFromContext(ctx)
		require.ErrorIs(t, err, ErrNoUserIDInContext)
	})

	t.Run("return error if id is not string", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextAuthIDKey, 123)

		_, err := GetUserIDFromContext(ctx)
		require.ErrorIs(t, err, ErrNoUserIDInContext)
	})
}

func TestPutUserIDIntoContext(t *testing.T) {
	t.Run("add user if into context", func(t *testing.T) {
		ctx := context.Background()
		ctx = PutUserIDIntoContext(ctx, "123")

		userIDFromContext, err := GetUserIDFromContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "123", userIDFromContext)
	})
}
