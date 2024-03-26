package main

import (
	"fmt"

	"github.com/vvjke314/mkc-backend/internal/app"
)

//	@title			MKC API
//	@version		1.0
//	@description	MK CLOUD backend service.
//	@contact.email	mail-bla-bla
//	@host			localhost:8080
//	@BasePath		/

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @schemes					http
func main() {
	app := app.NewApplication()
	err := app.Init()
	if err != nil {
		// switch to logger here
		fmt.Printf("[app.Init]: Can't initialize application: %s", err)
	}
	err = app.Run()
	if err != nil {
		// switch to logger here
		fmt.Printf("[app.Run] Error occured: %s", err)
	}
}
