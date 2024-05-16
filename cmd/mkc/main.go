package main

import (
	"fmt"

	"github.com/vvjke314/mkc-backend/internal/app"
)

type Application interface {
	Init() error
	Run() error
	Log(string)
}

// @title			MKC API
// @version		1.0
// @description	MK CLOUD backend service.
// @contact.email	mail@dump
// @host			localhost:8080
// @BasePath		/
// @schemes http
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	// Поменять на app.NewApplication()
	app := app.NewApplication()
	err := app.Init()
	if err != nil {
		// switch to logger here
		app.Log(fmt.Sprintf("[app.Init]: Can't initialize application: %s\n", err))
	}
	err = app.Run()
	if err != nil {
		// switch to logger here
		app.Log(fmt.Sprintf("[app.Run] Error occured: %s\n", err))
	}
}
