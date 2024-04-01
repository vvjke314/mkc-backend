package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vvjke314/mkc-backend/internal/pkg/db"
)

type Application struct {
	ctx  context.Context
	repo *db.Repo
}

func NewApplication() *Application {
	return &Application{}
}

// Init инициализирует сервис
func (app *Application) Init() error {
	app.ctx = context.Background()
	app.repo = db.NewRepo()
	err := app.repo.Init()
	if err != nil {
		return fmt.Errorf("[db.Init]: Can't initialize to database: %w", err)
	}
	return nil
}

// Run запускает сервис
func (app *Application) Run() error {
	// Подключение к бд
	err := app.repo.Connect()
	if err != nil {
		return fmt.Errorf("[db.Connect]: Can't connect to database: %w", err)
	}
	defer app.repo.Close()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// authorize
	r.POST("/login", Login)
	r.GET("/logout", Logout)
	r.POST("/signup", app.Signup)

	r.Run()

	return nil
}
