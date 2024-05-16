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
	app := app.NewApplication()
	err := app.Init()
	if err != nil {
		app.Log(fmt.Sprintf("[app.Init]: can't initialize application: %s\n", err), "service")
	}
	err = app.Run()
	if err != nil {
		app.Log(fmt.Sprintf("[app.Run] error occured: %s\n", err), "service")
	}
}
