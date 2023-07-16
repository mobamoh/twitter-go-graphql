package jwt

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	twitter_go_graphql "github.com/mobamoh/twitter-go-graphql"
	"github.com/mobamoh/twitter-go-graphql/config"
	"net/http"
	"time"
)

var SignatureType = jwa.HS256

type TokenService struct {
	Config *config.Config
}

func (ts *TokenService) ParseTokenFromRequest(ctx context.Context, req *http.Request) (twitter_go_graphql.AuthToken, error) {
	token, err := jwt.ParseRequest(req,
		jwt.WithValidate(true),
		jwt.WithIssuer(ts.Config.JWT.Issuer),
		jwt.WithVerify(SignatureType, []byte(ts.Config.JWT.Secret)),
	)
	if err != nil {
		return twitter_go_graphql.AuthToken{}, twitter_go_graphql.ErrInvalidAccessToken
	}
	return buildToken(token), nil
}

func buildToken(token jwt.Token) twitter_go_graphql.AuthToken {
	return twitter_go_graphql.AuthToken{
		ID:  token.JwtID(),
		Sub: token.Subject(),
	}
}

func (ts *TokenService) ParseToken(ctx context.Context, payload string) (twitter_go_graphql.AuthToken, error) {
	token, err := jwt.Parse([]byte(payload),
		jwt.WithValidate(true),
		jwt.WithIssuer(ts.Config.JWT.Issuer),
		jwt.WithVerify(SignatureType, []byte(ts.Config.JWT.Secret)),
	)
	if err != nil {
		return twitter_go_graphql.AuthToken{}, twitter_go_graphql.ErrInvalidAccessToken
	}
	return buildToken(token), nil
}

func (ts *TokenService) CreateRefreshToken(ctx context.Context, user twitter_go_graphql.User, tokenId string) (string, error) {
	token := jwt.New()
	if err := setDefaultToken(token, user, twitter_go_graphql.RefreshTokenLifetime, ts.Config); err != nil {
		return "", err
	}
	if err := token.Set(jwt.JwtIDKey, tokenId); err != nil {
		return "", fmt.Errorf("error set jwt ID key: %v", err)
	}
	signedToken, err := jwt.Sign(token, SignatureType, []byte(ts.Config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("error sign token: %v", err)
	}
	return string(signedToken), nil
}

func (ts *TokenService) CreateAccessToken(ctx context.Context, user twitter_go_graphql.User) (string, error) {
	token := jwt.New()
	if err := setDefaultToken(token, user, twitter_go_graphql.AccessTokenLifetime, ts.Config); err != nil {
		return "", err
	}
	signedToken, err := jwt.Sign(token, SignatureType, []byte(ts.Config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("error sign token: %v", err)
	}
	return string(signedToken), nil
}
func setDefaultToken(token jwt.Token, user twitter_go_graphql.User, lifetime time.Duration, conf *config.Config) error {
	if err := token.Set(jwt.SubjectKey, user.ID); err != nil {
		return fmt.Errorf("error set jwt sub: %v", err)
	}
	if err := token.Set(jwt.IssuerKey, conf.JWT.Issuer); err != nil {
		return fmt.Errorf("error set jwt issuer key: %v", err)
	}
	if err := token.Set(jwt.IssuedAtKey, time.Now().Unix()); err != nil {
		return fmt.Errorf("error set jwt issued at key: %v", err)
	}
	if err := token.Set(jwt.ExpirationKey, time.Now().Add(lifetime).Unix()); err != nil {
		return fmt.Errorf("error set jwt expiration key: %v", err)
	}
	return nil
}
