package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vvjke314/mkc-backend/internal/pkg/db"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
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

	customer := ds.Customer{
		Id:         uuid.New(),
		FirstName:  "Vladimir",
		SecondName: "Abramov",
		Login:      "vvjkee",
		Password:   "bufybuff2002",
		Email:      "vvjkee@mail.ru",
		Type:       0,
	}
	err = app.r.SignUpCustomer(customer)
	if err != nil {
		return fmt.Errorf("[db.SignUpCustomer]: Can't signup customer: %w", err)
	}

	administrator := ds.Administrator{
		Id:       uuid.New(),
		Name:     "Polina",
		Email:    "polina.andronova@mail.ru",
		Password: "lyblyuVovu",
	}
	err = app.r.SignUpAdministrator(administrator)
	if err != nil {
		return fmt.Errorf("[db.SignUpAdministrator]: Can't signup admin: %w", err)
	}

	return nil
}
