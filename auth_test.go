package twitter_go_graphql

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegisterInput_Sanitize(t *testing.T) {
	input := RegisterInput{
		UserName:        " Mo   ",
		Email:           "  mo@MAIL.com    ",
		Password:        "password",
		ConfirmPassword: "password",
	}
	expected := RegisterInput{
		UserName:        "Mo",
		Email:           "mo@mail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}
	input.Sanitize()
	assert.Equal(t, expected, input)
}
func TestRegisterInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input RegisterInput
		err   error
	}{
		{
			name: "valid",
			input: RegisterInput{
				UserName:        "Mo",
				Email:           "mo@mail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: nil,
		},
		{
			name: "invalid username",
			input: RegisterInput{
				UserName:        "M",
				Email:           "mo@mail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		},
		{
			name: "invalid email",
			input: RegisterInput{
				UserName:        "Mo",
				Email:           "mo",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		},
		{
			name: "invalid password",
			input: RegisterInput{
				UserName:        "Mo",
				Email:           "mo@mail.com",
				Password:        "pass",
				ConfirmPassword: "pass",
			},
			err: ErrValidation,
		},
		{
			name: "invalid password confirmation",
			input: RegisterInput{
				UserName:        "Mo",
				Email:           "mo@mail.com",
				Password:        "password",
				ConfirmPassword: "wrong",
			},
			err: ErrValidation,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.input.Validate()
			if testCase.err != nil {
				require.ErrorIs(t, err, testCase.err)
			} else {
				require.NoError(t, err, testCase.err)
			}
		})
	}
}

func TestLoginInput_Sanitize(t *testing.T) {
	input := LoginInput{
		Email:    "  mo@MAIL.com    ",
		Password: "password",
	}
	expected := LoginInput{
		Email:    "mo@mail.com",
		Password: "password",
	}
	input.Sanitize()
	assert.Equal(t, expected, input)
}

func TestLoginInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input LoginInput
		err   error
	}{
		{
			name: "valid",
			input: LoginInput{
				Email:    "mo@mail.com",
				Password: "password",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: LoginInput{
				Email:    "mo",
				Password: "password",
			},
			err: ErrValidation,
		},
		{
			name: "invalid password",
			input: LoginInput{
				Email:    "mo@mail.com",
				Password: "",
			},
			err: ErrValidation,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.input.Validate()
			if testCase.err != nil {
				require.ErrorIs(t, err, testCase.err)
			} else {
				require.NoError(t, err, testCase.err)
			}
		})
	}
}
