package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/dsn"
)

func DBConnect(ctx context.Context) (*pgx.Conn, error) {
	url, err := dsn.GetDSN()
	if err != nil {
		return nil, fmt.Errorf("[dsn.GetDSN]: Can't get data string name: %w", err)
	}

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("[pgx.Connect]: Can't establish connection: %w", err)
	}
	return conn, nil
}
