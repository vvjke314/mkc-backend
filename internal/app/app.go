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
	logFile, err := os.OpenFile(
		"application.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		return fmt.Errorf("[app.Init] %w", err)
	}
	app.logger = zerolog.New(logFile).With().Timestamp().Logger()
	app.repo = db.NewRepo()
	err = app.repo.Init()
	if err != nil {
		return fmt.Errorf("[db.Init] %w", err)
	}
	return nil
}

func (app *Application) Log(message, customerId string) {
	msg := fmt.Sprintf("$s:error:$s", customerId, message)
	app.logger.Error().Msg(msg)
}

func (app *Application) SuccessLog(message, customerId string) {
	msg := fmt.Sprintf("$s:success_request:$s", customerId, message)
	app.logger.Log().Msg(msg)
}

// Run запускает сервис
func (app *Application) Run() error {
	// Подключение к бд
	err := app.repo.Connect()
	if err != nil {
		return fmt.Errorf("[db.Connect]: can't connect to database: %w", err)
	}
	defer app.repo.Close()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// authorize
	r.POST("/login", app.Login)
	r.POST("/signup", app.Signup)

	/*
		administrator := r.Group("/admin")
		administrator.Use(AdminMiddleware())
		administrator.GET("/admin/unattached", app.GetAllUnattachedProjects)
		administrator.GET("/admin/attached", app.GetAllAttachedProjects)
		administrator.POST("/admin/attach/:project_id", app.AttachAdmin)
		administrator.POST("/admin/:project_id/discussion", app.SumbitEmail)
	*/

	authorized := r.Group("/")

	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/logout", app.Logout)

		// project
		authorized.GET("/projects", app.GetProjects)
		//authorized.GET("/project/:project_id", app.GetProjectInfo)
		authorized.POST("/project", app.CreateProject)
		authorized.Use(app.FullAccessControl()).PUT("/project/:project_id", app.UpdateProjectName)
		authorized.Use(app.FullAccessControl()).DELETE("/project/:project_id", app.DeleteProject)

		// file
		authorized.Use(app.FullAccessControl()).POST("/project/:project_id/file", app.UploadFile)
		authorized.Use(app.FullAccessControl()).POST("/project/:project_id/files", app.UploadFiles)
		authorized.Use(app.FullAccessControl()).DELETE("/project/:project_id/file", app.DeleteFile)
		authorized.Use(app.AccessControl()).GET("/project/:project_id/file/:file_id", app.DownloadFile)

		// note
		authorized.Use(app.FullAccessControl()).POST("/project/:project_id/note", app.CreateNote)
		//authorized.Use(app.FullAccessControl()).PUT("/project/:project_id/note/:note_id", app.UpdateNoteDeadline)
		//authorized.Use(app.FullAccessControl()).DELETE("/project/:project_id/note/:note_id", app.DelteNote)

		// participants
		authorized.Use(app.FullAccessControl()).POST("/participants/:project_id", app.AddParticipant)
		authorized.Use(app.FullAccessControl()).PUT("/participants/:project_id", app.UpdateParticipantAccess)
		authorized.Use(app.FullAccessControl()).DELETE("/participants/:project_id", app.DeleteParticipant)
	}

	r.Run()

	return nil
}
