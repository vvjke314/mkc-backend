package app

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/vvjke314/mkc-backend/internal/pkg/config"
	"github.com/vvjke314/mkc-backend/internal/pkg/db"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (app Application) Run() error {
	err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("[config.GetConfig]: Can't read config: %w", err)
	}

	_, err = db.DBConnect(viper.GetString("DATABASE_USERNAME"), viper.GetString("DATABASE_NAME"), viper.GetString("DATABASE_PORT"), viper.GetString("DATABASE_PASSWORD"), context.Background())
	if err != nil {
		return fmt.Errorf("[db.DBConnect]: Can't connect to database: %w", err)
	}

	return nil
}
