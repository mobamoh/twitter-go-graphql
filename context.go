package twitter_go_graphql

import (
	"context"
)

type contextKey string

var (
	contextAuthIDKey contextKey = "currentUserId"
)

func GetUserIDFromContext(ctx context.Context) (string, error) {
	if ctx.Value(contextAuthIDKey) == nil {
		return "", ErrNoUserIDInContext
	}
	userId, ok := ctx.Value(contextAuthIDKey).(string)
	if !ok {
		return "", ErrNoUserIDInContext
	}

	return userId, nil
}

func PutUserIDIntoContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextAuthIDKey, id)
}
