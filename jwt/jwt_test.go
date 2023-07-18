package jwt

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/jwt"
	twitter_go_graphql "github.com/mobamoh/twitter-go-graphql"
	"github.com/mobamoh/twitter-go-graphql/config"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	conf         *config.Config
	tokenService *TokenService
)

func TestMain(m *testing.M) {
	config.LoadEnv(".env.test")
	conf = config.New()
	tokenService = NewTokenService(conf)
	os.Exit(m.Run())
}

func TestTokenService_CreateAccessToken(t *testing.T) {
	t.Run("can create access token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter_go_graphql.User{
			ID: "123",
		}
		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		now = func() time.Time {
			return time.Now()
		}

		parsedToken, err := jwt.Parse([]byte(token),
			jwt.WithValidate(true),
			jwt.WithVerify(SignatureType, []byte(conf.JWT.Secret)),
			jwt.WithIssuer(conf.JWT.Issuer))
		require.NoError(t, err)
		require.Equal(t, user.ID, parsedToken.Subject())
		require.Equal(t,
			time.Now().Add(twitter_go_graphql.AccessTokenLifetime).Unix(),
			parsedToken.Expiration().Unix())

	})
}

func TestTokenService_CreateRefreshToken(t *testing.T) {
	t.Run("can create refresh token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter_go_graphql.User{
			ID: "123",
		}
		token, err := tokenService.CreateRefreshToken(ctx, user, "456")
		require.NoError(t, err)

		now = func() time.Time {
			return time.Now()
		}

		parsedToken, err := jwt.Parse([]byte(token),
			jwt.WithValidate(true),
			jwt.WithVerify(SignatureType, []byte(conf.JWT.Secret)),
			jwt.WithIssuer(conf.JWT.Issuer))
		require.NoError(t, err)
		require.Equal(t, user.ID, parsedToken.Subject())
		require.Equal(t, "456", parsedToken.JwtID())
		require.Equal(t,
			time.Now().Add(twitter_go_graphql.RefreshTokenLifetime).Unix(),
			parsedToken.Expiration().Unix())

	})
}

func TestTokenService_ParseToken(t *testing.T) {
	t.Run("can parse a valid access token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter_go_graphql.User{
			ID: "123",
		}
		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)
		parsedToken, err := tokenService.ParseToken(ctx, token)
		require.NoError(t, err)
		require.Equal(t, user.ID, parsedToken.Sub)
	})

	t.Run("can parse a valid refresh token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter_go_graphql.User{
			ID: "123",
		}
		token, err := tokenService.CreateRefreshToken(ctx, user, "456")
		require.NoError(t, err)
		parsedToken, err := tokenService.ParseToken(ctx, token)
		require.NoError(t, err)
		require.Equal(t, user.ID, parsedToken.Sub)
		require.Equal(t, "456", parsedToken.ID)
	})

	t.Run("return error if invalid access token", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.ParseToken(ctx, "invalid token")
		require.ErrorIs(t, err, twitter_go_graphql.ErrInvalidAccessToken)
	})

	t.Run("return error if access token expired", func(t *testing.T) {
		ctx := context.Background()
		user := twitter_go_graphql.User{
			ID: "123",
		}
		now = func() time.Time {
			return time.Now().Add(-twitter_go_graphql.AccessTokenLifetime * 5)
		}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)
		_, err = tokenService.ParseToken(ctx, token)
		require.ErrorIs(t, err, twitter_go_graphql.ErrInvalidAccessToken)
	})
}

func TestTokenService_ParseTokenFromRequest(t *testing.T) {
	t.Run("can parse access token from request", func(t *testing.T) {
		ctx := context.Background()
		user := twitter_go_graphql.User{
			ID: "123",
		}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", token)

		requestToken, err := tokenService.ParseTokenFromRequest(ctx, req)
		require.NoError(t, err)
		require.Equal(t, user.ID, requestToken.Sub)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		requestToken, err = tokenService.ParseTokenFromRequest(ctx, req)
		require.NoError(t, err)
		require.Equal(t, user.ID, requestToken.Sub)
	})

	t.Run("fail with expired token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter_go_graphql.User{
			ID: "123",
		}
		now = func() time.Time {
			return time.Now().Add(-twitter_go_graphql.AccessTokenLifetime * 5)
		}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", token)

		_, err = tokenService.ParseTokenFromRequest(ctx, req)
		require.ErrorIs(t, err, twitter_go_graphql.ErrInvalidAccessToken)
	})

	t.Run("fail with expired token", func(t *testing.T) {
		ctx := context.Background()

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "invalid token")

		_, err := tokenService.ParseTokenFromRequest(ctx, req)
		require.ErrorIs(t, err, twitter_go_graphql.ErrInvalidAccessToken)
	})
}
func teardownTimeNow(t *testing.T) {
	t.Helper()
	now = func() time.Time {
		return time.Now()
	}
}
