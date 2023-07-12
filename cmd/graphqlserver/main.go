package main

import (
	"context"
	"fmt"
	"github.com/mobamoh/twitter-go-graphql/config"
	"github.com/mobamoh/twitter-go-graphql/postgres"
	"log"
)

func main() {
	ctx := context.Background()
	config := config.New()
	db := postgres.New(ctx, config)
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("WORKING")
}
