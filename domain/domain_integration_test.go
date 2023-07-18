//go:build integration
// +build integration

package domain

import (
	"context"
	twitter "github.com/mobamoh/twitter-go-graphql"
	myConf "github.com/mobamoh/twitter-go-graphql/config"
	"github.com/mobamoh/twitter-go-graphql/jwt"
	"github.com/mobamoh/twitter-go-graphql/postgres"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"testing"
)

var (
	config           *myConf.Config
	db               *postgres.DB
	authTokenService twitter.AuthTokenService
	authService      twitter.AuthService
	userRepo         twitter.UserRepo
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	myConf.LoadEnv(".env.test")
	PasswordHashCost = bcrypt.MinCost

	config = myConf.New()
	db = postgres.New(ctx, config)
	defer db.Close()

	if err := db.Drop(); err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	userRepo = postgres.NewUserRepo(db)
	authTokenService = jwt.NewTokenService(config)
	authService = NewAuthService(userRepo, authTokenService)

	os.Exit(m.Run())
}
