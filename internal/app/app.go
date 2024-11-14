package app

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/vvjke314/mkc-backend/internal/pkg/config"
	"github.com/vvjke314/mkc-backend/internal/pkg/db"
)

type Application struct {
	ctx    context.Context
	repo   *db.Repo
	logger zerolog.Logger
	redis  *redis.Client
}

func NewApplication() *Application {
	return &Application{}
}

// Init инициализирует сервис
func (app *Application) Init() error {
	app.ctx = context.Background()
	logFile, err := os.OpenFile(
		"logs/application.log",
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
	err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("[config.GetConfig] %w", err)
	}
	adr := fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOSTNAME"), viper.GetString("REDIS_PORT"))
	app.redis = redis.NewClient(&redis.Options{
		Addr:     adr,
		Password: "",
		DB:       0,
	})

	str, err := app.redis.Ping(app.ctx).Result()
	fmt.Println(str)
	if err != nil {
		return fmt.Errorf("[redis.Init] %w", err)
	}
	return nil
}

func (app *Application) Log(message, customerId string) {
	msg := fmt.Sprintf("%s:error:%s", customerId, message)
	app.logger.Error().Msg(msg)
}

func (app *Application) SuccessLog(message, customerId string) {
	msg := fmt.Sprintf("%s:success_request:%s", customerId, message)
	app.logger.Log().Msg(msg)
}

// Run запускает сервис
func (app *Application) Run() error {
	// Подключение к бд
	err := app.repo.Connect()
	if err != nil {
		return fmt.Errorf("[repo.Connect]: can't connect to database: %w", err)
	}
	defer app.repo.Close()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// authorize
	r.POST("/login", app.Login)   // +
	r.POST("/signup", app.Signup) // +

	// administrator
	administrator := r.Group("/admin")
	administrator.POST("/signup", app.SignUpAdmin)                 // +-
	administrator.Use(app.BasicAuthMiddleware())                   //
	administrator.GET("/unattached", app.GetAllUnattachedProjects) // +
	administrator.GET("/attached", app.GetAllAttachedProjects)     // +-
	administrator.GET("/attach/:project_id", app.AttachAdmin)      // +
	administrator.POST("/:project_id/send", app.GetCustomerEmail)  // +-

	r.GET("/subscription/:customer_id", app.GetSubscription) //

	authorized := r.Group("/")

	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/logout", app.Logout) // +

		// subscription
		authorized.GET("/payment_info", app.GetPaymentUrl) //

		// project
		authorized.GET("/projects", app.GetProjects)                                           // +
		authorized.POST("/project", app.CheckSubscription(), app.CreateProject)                // +
		authorized.GET("/project/:project_id", app.AccessControl(), app.GetProjectInfo)        // +
		authorized.PUT("/project/:project_id", app.FullAccessControl(), app.UpdateProjectName) // +
		authorized.DELETE("/project/:project_id", app.FullAccessControl(), app.DeleteProject)  // +

		// file
		authorized.POST("/project/:project_id/file", app.FullAccessControl(), app.CheckSubscription(), app.UploadFile)   // +
		authorized.POST("/project/:project_id/files", app.FullAccessControl(), app.CheckSubscription(), app.UploadFiles) // +
		authorized.DELETE("/project/:project_id/file", app.FullAccessControl(), app.DeleteFile)                          // +
		authorized.GET("/project/:project_id/file/:file_id", app.AccessControl(), app.DownloadFile)                      // +
		authorized.GET("/project/:project_id/files", app.AccessControl(), app.GetFiles)                                  // +

		// note
		authorized.POST("/project/:project_id/note", app.FullAccessControl(), app.CreateNote)                 // +
		authorized.PUT("/project/:project_id/note/:note_id", app.FullAccessControl(), app.UpdateNoteDeadline) // +-
		authorized.DELETE("/project/:project_id/note/:note_id", app.FullAccessControl(), app.DeleteNote)      // +

		// participants
		authorized.POST("/participants/:project_id", app.FullAccessControl(), app.AddParticipant)         // +
		authorized.PUT("/participants/:project_id", app.FullAccessControl(), app.UpdateParticipantAccess) // +-
		authorized.DELETE("/participants/:project_id", app.FullAccessControl(), app.DeleteParticipant)    // +-
		authorized.GET("/participants/:project_id", app.AccessControl(), app.GetAllParticipants)          // +
	}

	r.Run()

	return nil
}
