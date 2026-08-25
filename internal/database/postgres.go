package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, databaseURL string) *pgxpool.Pool {
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal("failed to create database pool:", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		log.Fatal("failed to connect to database:", err)
	}

	log.Println("PostgreSQL connected successfully")
	return db
}
