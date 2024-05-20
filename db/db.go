package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPSQLStorage(connString string) (*pgxpool.Pool, error) {
	config, err := pgx.ParseConfig(connString)

	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), &pgxpool.Config{
		ConnConfig: &pgx.ConnConfig{
			Config: config.Config,
		},
	})

	return db, err
}
