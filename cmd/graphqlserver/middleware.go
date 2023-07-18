package main

import (
	twitter_go_graphql "github.com/mobamoh/twitter-go-graphql"
	"net/http"
)

func authMiddleware(service twitter_go_graphql.AuthTokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tokenFromRequest, err := service.ParseTokenFromRequest(ctx, r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx = twitter_go_graphql.PutUserIDIntoContext(ctx, tokenFromRequest.Sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
