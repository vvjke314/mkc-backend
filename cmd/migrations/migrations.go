package main

import (
	"flag"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	migrationsPath = "migrations"
	driver         = "postgres"
)

func main() {
	// Switch to logger here
	fmt.Println("Starting migrations")

	sql, err := goose.OpenDBWithDriver("pgx", "postgres://postgres:mkcdbpass@localhost:5432/mkcdb")
	if err != nil {
		// switch to logger here
		fmt.Println(err.Error())
	}

	isDownMigrate := flag.Bool("d", false, "do goose down migration")
	flag.Parse()

	// switch to logger here
	fmt.Println("Migrating")

	if *isDownMigrate == false {
		err = goose.Up(sql, "./migrations")
		if err != nil {
			// switch to logger here
			fmt.Println(err.Error())
		}
	} else {
		err = goose.DownTo(sql, "./migrations", 0)
		if err != nil {
			// switch to logger here
			fmt.Println(err.Error())
		}
	}
}
