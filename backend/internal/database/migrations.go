package database

import (
	"context"
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations executes all SQL migrations in order
func RunMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Create migrations tracking table
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Read migration files
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort migrations by filename (they should be numbered like 001_xxx.sql)
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	// Run each migration
	for _, filename := range files {
		version := strings.TrimSuffix(filename, ".sql")

		// Check if already applied
		var exists bool
		err := pool.QueryRow(ctx, 
			"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)",
			version,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if exists {
			continue
		}

		// Read and execute migration
		content, err := migrationsFS.ReadFile("migrations/" + filename)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", filename, err)
		}

		log.Printf("Running migration: %s", filename)

		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Record migration
		_, err = pool.Exec(ctx,
			"INSERT INTO schema_migrations (version) VALUES ($1)",
			version,
		)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		log.Printf("✓ Migration %s applied successfully", filename)
	}

	return nil
}
