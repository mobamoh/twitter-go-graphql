package twitter_go_graphql

import (
	"github.com/mobamoh/twitter-go-graphql/faker"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTweetInput_Sanitize(t *testing.T) {
	input := CreateTweetInput{
		Body: "   hello  ",
	}
	want := CreateTweetInput{
		Body: "hello",
	}

	input.Sanitize()
	require.Equal(t, want, input)
}

func TestCreateTweetInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input CreateTweetInput
		err   error
	}{
		{
			name: "valid",
			input: CreateTweetInput{
				Body: "hello",
			},
			err: nil,
		},
		{
			name: "tweet not long enough",
			input: CreateTweetInput{
				Body: "h",
			},
			err: ErrValidation,
		},
		{
			name: "tweet too long",
			input: CreateTweetInput{
				Body: faker.RandStringRunes(300),
			},
			err: ErrValidation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
