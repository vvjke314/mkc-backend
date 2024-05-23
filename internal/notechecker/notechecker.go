package notechecker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/vvjke314/mkc-backend/internal/pkg/config"
	"github.com/vvjke314/mkc-backend/internal/pkg/db"
)

type NoteChecker struct {
	ctx    context.Context
	repo   *db.Repo
	logger zerolog.Logger
}

func NewNoteChecker() *NoteChecker {
	return &NoteChecker{}
}

// Init инициализирует сервис проверки заметок
func (nc *NoteChecker) Init() error {
	nc.ctx = context.Background()
	logFile, err := os.OpenFile(
		"logs/notecheker.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		return fmt.Errorf("[noteChecker.Init][os.OpenFile]: %w", err)
	}
	nc.logger = zerolog.New(logFile).With().Timestamp().Logger()
	nc.repo = db.NewRepo()
	err = nc.repo.Init()
	if err != nil {
		return fmt.Errorf("[noteChecker.Init][repo.Init]: %w", err)
	}
	return nil
}

// Log логирует сообщения в указанный файл
func (nc *NoteChecker) Log(message string) {
	nc.logger.Error().Msg(message)
}

// SuccessLog логирует сообщения в указанный файл
func (nc *NoteChecker) SuccessLog(message string) {
	nc.logger.Log().Msg(message)
}

func (nc *NoteChecker) Run() error {
	// Подключение к БД
	err := nc.repo.Connect()
	if err != nil {
		return fmt.Errorf("repo.Connect: can't connect to database: %w", err)
	}
	// Закрываем подключение к БД
	defer nc.repo.Close()

	// Создание и запуск планировщика задач
	c := cron.New()
	// Задача на отправку уведомления за час до дедлайна
	err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("[config.GetConfig] can't parse config: %w", err)
	}
	timeCheck := fmt.Sprintf("*/%v * * * *", viper.GetInt("EMAIL_CHECK_TIME"))
	c.AddFunc(timeCheck, func() {
		err := nc.repo.ProccessNotes(time.Duration(viper.GetInt("EMAIL_BEFORE_TIME"))*time.Minute, "SELECT id, project_id, title, content, update_datetime, deadline, overdue FROM note WHERE deadline <= $1 AND overdue = 0")
		if err != nil {
			msg := fmt.Sprintf("Error notifying upcoming deadlines: %v", err)
			nc.Log(msg)
		}
		msg := fmt.Sprintf("[%vmin]checked", viper.GetInt("EMAIL_CHECK_TIME"))
		nc.SuccessLog(msg)
	})
	// Задача на отправку уведомления за день до дедлайна
	c.AddFunc("@daily", func() {
		err := nc.repo.ProccessNotes(24*time.Hour, "SELECT id, project_id, title, content, update_datetime, deadline, overdue FROM note WHERE deadline <= $1 AND overdue = 0")
		if err != nil {
			msg := fmt.Sprintf("Error notifying upcoming deadlines: %v", err)
			nc.Log(msg)
		}
		nc.SuccessLog("[daily]checked")
	})

	// Задача на отправку уведомления за час до дедлайна
	c.AddFunc("@hourly", func() {
		err := nc.repo.ProccessNotes(1*time.Hour, "SELECT id, project_id, title, content, update_datetime, deadline, overdue FROM note WHERE deadline <= $1 AND overdue = 0")
		if err != nil {
			msg := fmt.Sprintf("Error notifying upcoming deadlines: %v", err)
			nc.Log(msg)
		}
		nc.SuccessLog("[hourly]checked")
	})

	c.Start()

	// Ожидание завершения работы
	select {}
}
