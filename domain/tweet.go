package domain

import (
	"context"
	twitter "github.com/mobamoh/twitter-go-graphql"
	"github.com/mobamoh/twitter-go-graphql/uuid"
)

type TweetService struct {
	TweetRepo twitter.TweetRepo
}

func NewTweetService(repo twitter.TweetRepo) *TweetService {
	return &TweetService{
		TweetRepo: repo,
	}
}

func (t *TweetService) All(ctx context.Context) ([]twitter.Tweet, error) {
	return t.TweetRepo.All(ctx)
}

func (t *TweetService) Create(ctx context.Context, input twitter.CreateTweetInput) (twitter.Tweet, error) {
	currentUserID, err := twitter.GetUserIDFromContext(ctx)
	if err != nil {
		return twitter.Tweet{}, twitter.ErrUnauthenticated
	}
	input.Sanitize()
	if err := input.Validate(); err != nil {
		return twitter.Tweet{}, err
	}
	tweet, err := t.TweetRepo.Create(ctx, twitter.Tweet{
		Body:   input.Body,
		UserID: currentUserID,
	})
	if err != nil {
		return twitter.Tweet{}, err
	}
	return tweet, nil
}

func (t *TweetService) GetByID(ctx context.Context, id string) (twitter.Tweet, error) {
	if !uuid.Validate(id) {
		return twitter.Tweet{}, twitter.ErrInvalidUUID
	}
	return t.TweetRepo.GetByID(ctx, id)
}
