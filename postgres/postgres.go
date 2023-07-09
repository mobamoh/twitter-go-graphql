package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mobamoh/twitter-go-graphql/config"
	"log"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, config *config.Config) *DB {
	dbConf, err := pgxpool.ParseConfig(config.Database.URL)
	if err != nil {
		log.Fatalf("cannot parse postgres config %+v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, dbConf)
	if err != nil {
		log.Fatalf("cannot connect to postgres %+v", err)
	}
	db := &DB{pool}
	db.Ping(ctx)
	return db
}
func (db *DB) Ping(ctx context.Context) {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatalf("cannot ping postgres %+v", err)
	}
	log.Panicln("postgres connected")
}
func (db *DB) Close() {
	db.Pool.Close()
}
