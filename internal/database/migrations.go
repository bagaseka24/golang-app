package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func RunMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Create users table
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil
}