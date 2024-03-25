package app

import (
	"context"
	"fmt"

	"github.com/vvjke314/mkc-backend/internal/pkg/db"
	"github.com/vvjke314/mkc-backend/internal/pkg/service"
)

type Application struct {
	ctx context.Context
	r   *db.Repo
	srv *service.Service
}

func NewApplication() *Application {
	return &Application{}
}

// Init
// Инициализирует сервис
func (app *Application) Init() {
	app.ctx = context.Background()
	app.r = db.NewRepo()
	app.r.Init()
}

// Run
// Запускает сервис
func (app *Application) Run() error {
	// подключение к бд
	conn, err := app.r.Connect(app.ctx)
	if err != nil {
		return fmt.Errorf("[db.DBConnect]: Can't connect to database: %w", err)
	}
	defer conn.Close()

	return nil
}
