package main

import (
	"fmt"

	"github.com/pressly/goose/v3"

	"github.com/spf13/viper"
	"github.com/vvjke314/mkc-backend/internal/pkg/config"
)

const (
	migrationsPath = "migrations"
	driver         = "postgres"
)

func main() {
	// Switch to logger here
	fmt.Println("Starting migrations")

	config.GetConfig()

	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=%s sslmode=disable", viper.GetString("DATABASE_USERNAME"), viper.GetString("DATABASE_PASSWORD"), viper.GetString("DATABASE_NAME"), viper.GetString("DATABASE_PORT"))

	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		// Switch to logger here
		fmt.Println(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			// Switch to logger here
			fmt.Println(err.Error())
		}
	}()

	err = goose.Up(db, dsn)
	if err != nil {
		// Switch to logger here
		fmt.Println(err)
	}
}
