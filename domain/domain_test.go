package domain

import (
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	PasswordHashCost = bcrypt.MinCost
	os.Exit(t.Run())
}
