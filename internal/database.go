package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func ConnectDB(url string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), url)
}
