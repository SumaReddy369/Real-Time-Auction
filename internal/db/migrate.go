package db

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS


func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	if err := ensureMigrationsTable(ctx, pool); err != nil {
		return err
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		applied, err := isApplied(ctx, pool, f)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := migrationsFS.ReadFile("migrations/" + f)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}

		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("run migration %s: %w", f, err)
		}

		if _, err := pool.Exec(ctx, `INSERT INTO schema_migrations (name) VALUES ($1)`, f); err != nil {
			return fmt.Errorf("record migration %s: %w", f, err)
		}
	}
	return nil
}

func ensureMigrationsTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	return err
}

func isApplied(ctx context.Context, pool *pgxpool.Pool, name string) (bool, error) {
	var count int
	err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM schema_migrations WHERE name = $1`, name).Scan(&count)
	return count > 0, err
}
