//go:generate go run github.com/vektah/dataloaden UserLoader string *twitter-go-graphql/graph.User

package graph

import (
	"context"
	"fmt"
	twitter_go_graphql "github.com/mobamoh/twitter-go-graphql"
	"net/http"
	"time"
)

const loaderskey = "dataloaders"

type Loaders struct {
	UserByID UserLoader
}

type Repos struct {
	UserRepo twitter_go_graphql.UserRepo
}

func DataloaderMiddleware(repos *Repos) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loaderskey, &Loaders{
				UserByID: UserLoader{
					wait:     1 * time.Millisecond,
					maxBatch: 10,
					fetch: func(ids []string) ([]*User, []error) {
						users, err := repos.UserRepo.GetByIDs(r.Context(), ids)
						if err != nil {
							return nil, []error{err}
						}

						userMap := map[string]*User{}
						for _, user := range users {
							userMap[user.ID] = mapUser(user)
						}

						result := make([]*User, len(ids))
						for i, id := range ids {
							user, ok := userMap[id]
							if !ok {
								return nil, []error{fmt.Errorf("user with id %s is missing", id)}
							}
							result[i] = user
						}
						return result, nil
					},
				},
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func DataloaderFor(ctx context.Context) *Loaders {
	return ctx.Value(loaderskey).(*Loaders)
}
