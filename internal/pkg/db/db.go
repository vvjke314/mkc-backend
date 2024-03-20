package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func DBConnect(dbUserName, dbName, dbPort, dbPassword string, ctx context.Context) (*pgx.Conn, error) {
	url := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", dbUserName, dbPassword, dbPort, dbName)
	fmt.Println(url)
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
