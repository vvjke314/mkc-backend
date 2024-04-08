package app

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/vvjke314/mkc-backend/internal/pkg/db"
)

type Application struct {
	ctx    context.Context
	repo   *db.Repo
	logger zerolog.Logger
}

func NewApplication() *Application {
	return &Application{}
}

// Init инициализирует сервис
func (app *Application) Init() error {
	app.ctx = context.Background()
	file, err := os.OpenFile(
		"application.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		return fmt.Errorf("[app.Init] %w", err)
	}
	app.logger = zerolog.New(file).With().Timestamp().Logger()
	app.repo = db.NewRepo()
	err = app.repo.Init()
	if err != nil {
		return fmt.Errorf("[db.Init] %w", err)
	}
	return nil
}

func (app *Application) Log(message string) {
	app.logger.Error().Msg(message)
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
	r.POST("/login", app.Login)
	r.POST("/signup", app.Signup)

	authorized := r.Group("/")

	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/logout", app.Logout)

		//project
		authorized.POST("/project", app.CreateProject)
		authorized.PUT("/project/:project_id", app.UpdateProjectName)
	}
	//

	r.Run()

	return nil
}
