package app

import (
	"context"
	"fmt"

	"github.com/vvjke314/mkc-backend/internal/pkg/db"
)

type Application struct {
	ctx context.Context
	r   *db.Repo
	//srv *service.Service
}

func NewApplication() *Application {
	return &Application{}
}

// Init
// Инициализирует сервис
func (app *Application) Init() error {
	app.ctx = context.Background()
	app.r = db.NewRepo()
	err := app.r.Init()
	if err != nil {
		return fmt.Errorf("[db.Init]: Can't initialize to database: %w", err)
	}

	//TO-DO: SERVICE INIT
	return nil
}

// Run
// Запускает сервис
func (app *Application) Run() error {
	// подключение к бд
	err := app.r.Connect()
	if err != nil {
		return fmt.Errorf("[db.Connect]: Can't connect to database: %w", err)
	}
	defer app.r.Close()

	return nil
}
