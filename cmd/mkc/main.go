package main

import (
	"fmt"

	"github.com/vvjke314/mkc-backend/internal/app"
)

func main() {
	app := app.NewApplication()
	err := app.Run()
	if err != nil {
		fmt.Printf("[app.Run] Error occured: %s", err)
	}
}
