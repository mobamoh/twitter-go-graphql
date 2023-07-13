package postgres

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mobamoh/twitter-go-graphql/config"
	"log"
	"path"
	"runtime"
)

type DB struct {
	Pool   *pgxpool.Pool
	config *config.Config
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
	db := &DB{pool, config}
	db.Ping(ctx)
	return db
}
func (db *DB) Ping(ctx context.Context) {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatalf("cannot ping postgres %+v", err)
	}
	log.Println("postgres connected")
}
func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) Migrate() error {
	_, file, _, _ := runtime.Caller(0)
	migrationPath := fmt.Sprintf("file:///%s/migrations", path.Dir(file))
	m, err := migrate.New(migrationPath, db.config.Database.URL)
	if err != nil {
		return fmt.Errorf("error creating the migrate instance: %v", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error migrating up: %v", err)
	}
	fmt.Println("migration done")
	return nil
}

func (db *DB) Drop() error {
	_, file, _, _ := runtime.Caller(0)
	migrationPath := fmt.Sprintf("file:///%s/migrations", path.Dir(file))
	m, err := migrate.New(migrationPath, db.config.Database.URL)
	if err != nil {
		return fmt.Errorf("error creating the migrate instance: %v", err)
	}
	if err = m.Drop(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error dropping: %v", err)
	}
	fmt.Println("migration drop")
	return nil
}

func (db *DB) Truncate(ctx context.Context) error {
	if _, err := db.Pool.Exec(ctx, `DELETE FROM users`); err != nil {
		return fmt.Errorf("error truncating: %v", err)
	}
	return nil
}
