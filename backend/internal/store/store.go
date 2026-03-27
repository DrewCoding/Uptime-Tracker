package store

import (
	"context"
	"database/sql"
	"fmt"
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
