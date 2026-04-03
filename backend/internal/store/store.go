package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

// CheckRow represents a health check result read from the database.
type CheckRow struct {
	ID         int       `json:"id"`
	URL        string    `json:"url"`
	StatusCode *int      `json:"status_code"`
	LatencyMs  int64     `json:"latency_ms"`
	CheckedAt  time.Time `json:"checked_at"`
}

// GetLatestChecks returns the most recent health check for every monitored URL.
func (db *DB) GetLatestChecks(ctx context.Context) ([]CheckRow, error) {
	query := `
		SELECT DISTINCT ON (url) id, url, status_code, latency_ms, checked_at
		FROM health_checks
		ORDER BY url, checked_at DESC
	`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query latest checks: %w", err)
	}
	defer rows.Close()

	var results []CheckRow
	for rows.Next() {
		var r CheckRow
		if err := rows.Scan(&r.ID, &r.URL, &r.StatusCode, &r.LatencyMs, &r.CheckedAt); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

// GetCheckHistory returns the most recent `limit` check results for a specific URL.
func (db *DB) GetCheckHistory(ctx context.Context, url string, limit int) ([]CheckRow, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	query := `
		SELECT id, url, status_code, latency_ms, checked_at
		FROM health_checks
		WHERE url = $1
		ORDER BY checked_at DESC
		LIMIT $2
	`
	rows, err := db.conn.QueryContext(ctx, query, url, limit)
	if err != nil {
		return nil, fmt.Errorf("query check history: %w", err)
	}
	defer rows.Close()

	var results []CheckRow
	for rows.Next() {
		var r CheckRow
		if err := rows.Scan(&r.ID, &r.URL, &r.StatusCode, &r.LatencyMs, &r.CheckedAt); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}
