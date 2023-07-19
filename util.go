package twitter_go_graphql

import (
	"errors"
)

var (
	ErrValidation         = errors.New("validation error")
	ErrNotFound           = errors.New("not found error")
	ErrServer             = errors.New("server error")
	ErrBadCredentials     = errors.New("bad credentials")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrNoUserIDInContext  = errors.New("missing user id in context")
	ErrGenAccessToken     = errors.New("generate access token error")
	ErrUnauthenticated    = errors.New("unauthenticated error")
	ErrInvalidUUID        = errors.New("invalid uuid error")
	ErrForbidden          = errors.New("forbidden error")
)
