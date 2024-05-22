package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPSQLStorage(connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)

	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)

	return db, err
}
