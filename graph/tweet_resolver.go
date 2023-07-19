package graph

import (
	"context"
	twitter "github.com/mobamoh/twitter-go-graphql"
)

func mapTweet(tweet twitter.Tweet) *Tweet {
	return &Tweet{
		ID:        tweet.ID,
		Body:      tweet.Body,
		UserID:    tweet.UserID,
		CreatedAt: &tweet.CreatedAt,
	}
}

func MapTweets(tweets []twitter.Tweet) []*Tweet {
	allTweets := make([]*Tweet, len(tweets))
	for i, tweet := range tweets {
		allTweets[i] = mapTweet(tweet)
	}
	return allTweets
}

func (q *queryResolver) Tweets(ctx context.Context) ([]*Tweet, error) {
	all, err := q.TweetService.All(ctx)
	if err != nil {
		return nil, err
	}
	return MapTweets(all), nil
}

func (m mutationResolver) CreateTweet(ctx context.Context, input CreateTweetInput) (*Tweet, error) {
	tweet, err := m.TweetService.Create(ctx, twitter.CreateTweetInput{
		Body: input.Body,
	})
	if err != nil {
		return nil, buildError(ctx, err)
	}
	return mapTweet(tweet), nil
}

func (m mutationResolver) DeleteTweet(ctx context.Context, id string) (bool, error) {
	if err := m.TweetService.Delete(ctx, id); err != nil {
		return false, buildError(ctx, err)
	}
	return true, nil
}

func (t tweetResolver) User(ctx context.Context, obj *Tweet) (*User, error) {
	return DataloaderFor(ctx).UserByID.Load(obj.UserID)

	//user, err := t.UserService.GetByID(ctx, obj.UserID)
	//if err != nil {
	//	return nil, buildError(ctx, err)
	//}
	//return mapUser(user), nil
}

func (m mutationResolver) CreateReply(ctx context.Context, parentID string, input CreateTweetInput) (*Tweet, error) {
	reply, err := m.TweetService.CreateReply(ctx, parentID, twitter.CreateTweetInput{
		Body: input.Body,
	})
	if err != nil {
		return nil, buildError(ctx, err)
	}
	return mapTweet(reply), nil
}
