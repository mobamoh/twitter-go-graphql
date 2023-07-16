package twitter_go_graphql

import (
	"fmt"
)

var (
	ErrValidation     = fmt.Errorf("validation error")
	ErrNotFound       = fmt.Errorf("not found error")
	ErrServer         = fmt.Errorf("server error")
	ErrBadCredentials = fmt.Errorf("bad credentials")

	ErrInvalidAccessToken = fmt.Errorf("invalid access token")
)
