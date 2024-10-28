package pg

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	DB *pgxpool.Pool
}

func DBPoolProvider(dsn string, ctx context.Context) (*DB, error) {
	connConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Println("error parsing config", err)
		return nil, err
	}

	if db, err := pgxpool.ConnectConfig(ctx, connConfig); err != nil {
		log.Println("error connecting to database", err)
		return nil, err
	} else {
		return &DB{DB: db}, nil
	}
}
