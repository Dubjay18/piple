package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	dbq "github.com/valentineejk/piple/db/sqlc"
)

func Connection() (*dbq.Queries, *pgxpool.Pool) {
	DATABASE_URL := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		log.Fatalf("database ping failed: %v", err)
	}
	fmt.Println("connected to postgres")

	return dbq.New(pool), pool
}
