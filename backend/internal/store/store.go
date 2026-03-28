package store

import (
	"context"
	"database/sql"
	"fmt"

	"tracker/internal/monitor"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func New(host string, port int, user, password, dbname string) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if err := conn.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &DB{conn: conn}, nil
}

// SaveCheck inserts a single health check result into the database.
func (db *DB) SaveCheck(ctx context.Context, r monitor.HealthCheck) error {
	query := `
		INSERT INTO health_checks (url, status_code, latency_ms, checked_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.conn.ExecContext(ctx, query, r.URL, r.StatusCode, r.LatencyMs, r.CheckedAt)
	if err != nil {
		return fmt.Errorf("insert health_check: %w", err)
	}
	return nil
}

// SaveChecks inserts multiple health check results in a single transaction.
func (db *DB) SaveChecks(ctx context.Context, results []monitor.HealthCheck) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO health_checks (url, status_code, latency_ms, checked_at)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	for _, r := range results {
		if _, err := stmt.ExecContext(ctx, r.URL, r.StatusCode, r.LatencyMs, r.CheckedAt); err != nil {
			return fmt.Errorf("exec for %s: %w", r.URL, err)
		}
	}

	return tx.Commit()
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}
