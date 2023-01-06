package pg

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	pgURL := os.Getenv("PG_URL")
	fmt.Println(pgURL)
	return pgx.Connect(ctx, pgURL)
}
