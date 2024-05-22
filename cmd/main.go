package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/adhupraba/ecom/cmd/api"
	"github.com/adhupraba/ecom/config"
	"github.com/adhupraba/ecom/db"
)

func main() {
	pgx, err := db.NewPSQLStorage(config.Envs.DbUrl)

	if err != nil {
		log.Fatal(err)
	}

	initStorage(pgx)

	server := api.NewAPIServer("0.0.0.0:"+config.Envs.ServerPort, pgx)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *pgxpool.Pool) {
	err := db.Ping(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
