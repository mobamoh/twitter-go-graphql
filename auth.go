package twitter_go_graphql

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var (
	UsernameMinLength = 2
	PasswordMinLength = 6
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (AuthResponse, error)
}

type AuthToken struct {
	ID  string
	Sub string
}

type AuthResponse struct {
	AccessToken string
	User        User
}
type RegisterInput struct {
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
}

func (r *RegisterInput) Sanitize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Username = strings.TrimSpace(r.Username)
}

func (r RegisterInput) Validate() error {
	if len(r.Username) < UsernameMinLength {
		return fmt.Errorf("%w: username must be at least %d characters long", ErrValidation, UsernameMinLength)
	}

	if !emailRegex.MatchString(r.Email) {
		return fmt.Errorf("%w: email has wrong format", ErrValidation)
	}

	if len(r.Password) < PasswordMinLength {
		return fmt.Errorf("%w: password must be at least %d characters long", ErrValidation, PasswordMinLength)
	}

	if r.Password != r.ConfirmPassword {
		return fmt.Errorf("%w: password confimation doesn't match", ErrValidation)
	}
	return nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (r *LoginInput) Sanitize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
}

func (r LoginInput) Validate() error {
	if !emailRegex.MatchString(r.Email) {
		return fmt.Errorf("%w: email has wrong format", ErrValidation)
	}
	if len(r.Password) == 0 {
		return fmt.Errorf("%w: password not specified", ErrValidation)
	}
	return nil
}
