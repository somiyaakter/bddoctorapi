package main

import (
	"context"
	"datalab_api/internal/database"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const migrationsDir = "migrations"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down>")
	}
	direction := os.Args[1]
	if direction != "up" && direction != "down" {
		log.Fatalf("unknown command %q, expected 'up' or 'down'", direction)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	ctx := context.Background()
	db := database.NewPostgresDB(ctx, databaseURL)
	defer db.Close()

	if err := ensureMigrationsTable(ctx, db); err != nil {
		log.Fatalf("failed to prepare schema_migrations table: %v", err)
	}

	var err error
	if direction == "up" {
		err = migrateUp(ctx, db)
	} else {
		err = migrateDown(ctx, db)
	}
	if err != nil {
		log.Fatalf("migration %s failed: %v", direction, err)
	}
}

func ensureMigrationsTable(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	return err
}

// migrationFile pairs a version's up/down SQL files, e.g. version
// "0001_create_doctors_and_chambers" with its .up.sql and .down.sql.
type migrationFile struct {
	version  string
	upPath   string
	downPath string
}

func loadMigrations() ([]migrationFile, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("reading migrations dir: %w", err)
	}

	byVersion := map[string]*migrationFile{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()

		var version, kind string
		switch {
		case strings.HasSuffix(name, ".up.sql"):
			version, kind = strings.TrimSuffix(name, ".up.sql"), "up"
		case strings.HasSuffix(name, ".down.sql"):
			version, kind = strings.TrimSuffix(name, ".down.sql"), "down"
		default:
			continue
		}

		m, ok := byVersion[version]
		if !ok {
			m = &migrationFile{version: version}
			byVersion[version] = m
		}
		path := filepath.Join(migrationsDir, name)
		if kind == "up" {
			m.upPath = path
		} else {
			m.downPath = path
		}
	}

	var result []migrationFile
	for _, m := range byVersion {
		result = append(result, *m)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].version < result[j].version })
	return result, nil
}

func appliedVersions(ctx context.Context, db *pgxpool.Pool) (map[string]bool, error) {
	rows, err := db.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := map[string]bool{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		applied[v] = true
	}
	return applied, rows.Err()
}

func migrateUp(ctx context.Context, db *pgxpool.Pool) error {
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}
	applied, err := appliedVersions(ctx, db)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if applied[m.version] {
			log.Printf("skipping already-applied migration: %s", m.version)
			continue
		}
		if m.upPath == "" {
			return fmt.Errorf("migration %s has no .up.sql file", m.version)
		}

		sqlBytes, err := os.ReadFile(m.upPath)
		if err != nil {
			return fmt.Errorf("reading %s: %w", m.upPath, err)
		}

		tx, err := db.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("applying %s: %w", m.version, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, m.version); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("recording %s: %w", m.version, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("committing %s: %w", m.version, err)
		}

		log.Printf("applied migration: %s", m.version)
	}
	return nil
}

func migrateDown(ctx context.Context, db *pgxpool.Pool) error {
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}
	applied, err := appliedVersions(ctx, db)
	if err != nil {
		return err
	}

	var latest *migrationFile
	for i := len(migrations) - 1; i >= 0; i-- {
		if applied[migrations[i].version] {
			latest = &migrations[i]
			break
		}
	}
	if latest == nil {
		log.Println("no applied migrations to roll back")
		return nil
	}
	if latest.downPath == "" {
		return fmt.Errorf("migration %s has no .down.sql file", latest.version)
	}

	sqlBytes, err := os.ReadFile(latest.downPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", latest.downPath, err)
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("reverting %s: %w", latest.version, err)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM schema_migrations WHERE version = $1`, latest.version); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("removing migration record %s: %w", latest.version, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing rollback of %s: %w", latest.version, err)
	}

	log.Printf("rolled back migration: %s", latest.version)
	return nil
}
