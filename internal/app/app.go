package app

import (
	"context"
	"fmt"

	"github.com/vvjke314/mkc-backend/internal/pkg/db"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (app Application) Run() error {
	ctx := context.Background()
	conn, err := db.DBConnect(ctx)
	if err != nil {
		return fmt.Errorf("[db.DBConnect]: Can't connect to database: %w", err)
	}
	defer conn.Close(ctx)

	return nil
}
