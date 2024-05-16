package notechecker

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
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
		"notecheker.log",
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

func (nc *NoteChecker) Run() error {
	// Подключение к БД
	err := nc.repo.Connect()
	if err != nil {
		return fmt.Errorf("repo.Connect: can't connect to database: %w", err)
	}
	// Закрываем подключение к БД
	defer nc.repo.Close()

	// СЮДА НАШУ РАБОТУ ДОБАВЛЯЕМ

	return nil
}
