package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	twitter_go_graphql "github.com/mobamoh/twitter-go-graphql"
)

type TweetRepo struct {
	db *DB
}

func NewTweetRepo(db *DB) *TweetRepo {
	return &TweetRepo{
		db: db,
	}
}

func (t *TweetRepo) All(ctx context.Context) ([]twitter_go_graphql.Tweet, error) {
	return listTweets(ctx, t.db.Pool)
}

func (t *TweetRepo) Create(ctx context.Context, tweet twitter_go_graphql.Tweet) (twitter_go_graphql.Tweet, error) {
	tx, err := t.db.Pool.Begin(ctx)
	if err != nil {
		return twitter_go_graphql.Tweet{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	tweet, err = createTweet(ctx, tx, tweet)
	if err != nil {
		return twitter_go_graphql.Tweet{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return twitter_go_graphql.Tweet{}, fmt.Errorf("error commiting: %v", err)
	}
	return tweet, nil
}

func (t *TweetRepo) GetByID(ctx context.Context, id string) (twitter_go_graphql.Tweet, error) {
	return getTweetByID(ctx, t.db.Pool, id)
}

func (t *TweetRepo) Delete(ctx context.Context, id string) error {
	tx, err := t.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	if err := deleteTweet(ctx, tx, id); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error commiting: %v", err)
	}
	return nil
}

func createTweet(ctx context.Context, tx pgx.Tx, tweet twitter_go_graphql.Tweet) (twitter_go_graphql.Tweet, error) {
	query := `INSERT INTO tweets(body, user_id,parent_id) VALUES ($1,$2,$3) RETURNING *;`
	newTweet := twitter_go_graphql.Tweet{}
	if err := pgxscan.Get(ctx, tx, &newTweet, query, tweet.Body, tweet.UserID, tweet.ParentID); err != nil {
		return twitter_go_graphql.Tweet{}, fmt.Errorf("error insert: %v", err)
	}
	return newTweet, nil
}

func getTweetByID(ctx context.Context, querier pgxscan.Querier, id string) (twitter_go_graphql.Tweet, error) {
	query := `SELECT * FROM tweets WHERE id = $1 LIMIT 1;`
	tweet := twitter_go_graphql.Tweet{}
	if err := pgxscan.Get(ctx, querier, &tweet, query, id); err != nil {
		if pgxscan.NotFound(err) {
			return twitter_go_graphql.Tweet{}, twitter_go_graphql.ErrNotFound
		}
		return twitter_go_graphql.Tweet{}, fmt.Errorf("error getting tweet: %+v", err)
	}
	return tweet, nil
}

func listTweets(ctx context.Context, querier pgxscan.Querier) ([]twitter_go_graphql.Tweet, error) {
	query := `SELECT * FROM tweets ORDER BY created_at DESC;`
	var tweets []twitter_go_graphql.Tweet
	if err := pgxscan.Select(ctx, querier, &tweets, query); err != nil {
		return nil, fmt.Errorf("error list tweets: %+v", err)
	}
	return tweets, nil
}

func deleteTweet(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM tweets WHERE id = $1;`
	if _, err := tx.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("error deleting: %v", err)
	}
	return nil
}
